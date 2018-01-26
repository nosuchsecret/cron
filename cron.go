package cron

import (
	//"log"
	"runtime"
	//"sync"
	"fmt"
	"time"
)

type Cron struct {
	stop     chan struct{}
	control  chan *Entry

	running  bool
	//log      log.Log

	entries  map[string]*Entry
	//lock     sync.Mutex

	jtree    *Rbtree
}

type Job interface {
	Run(*JobMeta, time.Time)
}

// The Schedule describes a job's duty cycle.
type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	// The schedule on which this job should be run.
	Op       int
	Schedule Schedule
	Meta     *JobMeta
	Node     *RbtreeNode

	// The next time the job will run. This is the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next     time.Time

	// The last time this job was run. This is the zero time if the job has never
	// been run.
	Prev     time.Time

	// The Job to run.
	Job      Job
}

// New returns a new Cron job runner, in the Local time zone.
func New() *Cron {
	return &Cron{
		control: make(chan *Entry),
		stop:    make(chan struct{}),
		entries: make(map[string]*Entry),
		jtree:   RbtreeInit(CronInsert),
		//lock
		running: false,
	}
}

// A wrapper that turns a func() into a cron.Job
type FuncJob func(m *JobMeta, next time.Time)

func (f FuncJob) Run(m *JobMeta, next time.Time) { f(m, next) }

func (c *Cron) DeleteFunc(jobid string) {
	c.DeleteJob(jobid)
}

func (c *Cron) DeleteJob(jobid string) {
	var entry *Entry

	if c.running {
		entry = &Entry{
			Op: 1,
			Meta: &JobMeta{Id: jobid},
		}
		c.control <- entry
		return
	}

	if entry = c.entries[jobid]; entry != nil {
		c.jtree.RbtreeDelete(entry.Node)
		c.entries[jobid] = nil
	}
}

func (c *Cron) AddFunc(spec string, meta *JobMeta, cmd func(*JobMeta, time.Time)) error {
	return c.AddJob(spec, meta, FuncJob(cmd))
}

func (c *Cron) AddJob(spec string, meta *JobMeta, cmd Job) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	c.Schedule(schedule, meta, cmd)
	return nil
}

func (c *Cron) Schedule(schedule Schedule, meta *JobMeta, cmd Job) {
	entry := &Entry{
		Schedule: schedule,
		Job:      cmd,
		Meta:     meta,
		Node:     &RbtreeNode{},
		Next:     schedule.Next(time.Now()),
	}
	entry.Node.data = entry

	if c.running {
		entry.Op = 0
		c.control <- entry
		return
	}
	c.entries[meta.Id] = entry
	c.jtree.RbtreeInsert(entry.Node)

	return
}

func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Cron) runWithRecovery(meta *JobMeta, j Job, next time.Time) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("cron: panic running job: %v\n%s\n", r, buf)
		}
	}()
	j.Run(meta, next)
}

func (c *Cron) run() {
	var meta *JobMeta
	var effective time.Time
	var node *RbtreeNode
	var entry *Entry
	now := time.Now()

	for {
		if c.jtree.NodeNum() == 0 {
			effective = now.AddDate(10, 0, 0)
		} else {
			node = c.jtree.FindMin()
			if e, ok := node.data.(*Entry); ok {
				entry = e
			}
			//TODO: maybe next is passed!
			effective = time.Time(entry.Next)
		}

		timer := time.NewTimer(effective.Sub(now))
		select {
		case now = <-timer.C:
			for {
				node = c.jtree.FindMin()
				if node == nil {
					break
				}
				if e, ok := node.data.(*Entry); ok {
					entry = e
				} //should not reach else

				if entry.Next.Before(now) {
					c.jtree.RbtreeDelete(node)

					next := entry.Schedule.Next(now)
					go c.runWithRecovery(entry.Meta, entry.Job, next)

					entry.Prev = entry.Next
					//entry.Next = entry.Schedule.Next(now)
					entry.Next = next

					c.jtree.RbtreeInsert(node)
				} else {
					break
				}
			}

		case entry := <-c.control:
			if entry.Op == 0 { //add
				meta = entry.Meta
				c.entries[meta.Id] = entry
				c.jtree.RbtreeInsert(entry.Node)
			} else if entry.Op == 1 { //delete
				jobid := entry.Meta.Id
				if entry = c.entries[jobid]; entry != nil {
					c.jtree.RbtreeDelete(entry.Node)
					c.entries[jobid] = nil
				}
			}

		case <-c.stop:
			timer.Stop()
			return
		}

		now = time.Now()
		timer.Stop()
	}
}

func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func CronInsert(temp, node, sentinel *RbtreeNode) {
	var p **RbtreeNode
	var ne, te *Entry
	for {
		if node.key < temp.key {
			fmt.Println("set p to left")
			p = &temp.left
		} else if node.key > temp.key {
			fmt.Println("set p to right")
			p = &temp.right
		} else {
			if e, ok := node.data.(*Entry); ok {
				ne = e
			}
			if e, ok := temp.data.(*Entry); ok {
				te = e
			}

			if ne.Next.Before(te.Next) {
				p = &temp.left
			} else {
				p = &temp.right
			}
		}

		if *p == sentinel {
			break
		}

		temp = *p
	}

	*p = node
	node.parent = temp
	node.left = sentinel
	node.right = sentinel
	rbtRed(node)
}
