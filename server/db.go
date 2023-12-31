package main

import (
  "log"
  "database/sql"
  "time"
  _ "github.com/mattn/go-sqlite3"
)

type Connection struct {
  filepath string
}
// ok this saves me literally only like 3 lines of code... but...
func (conn Connection) Do(fn func(db *sql.DB)) {
  // the boiler plate
  db,err := sql.Open("sqlite3",conn.filepath)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // the actual thing
  fn(db)
}

func (conn Connection) Init() bool {
  db,err := sql.Open("sqlite3",conn.filepath)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  _,err = db.Exec(`
    create table images(
      id integer primary key autoincrement, 
      path varchar, 
      name varchar, 
      date datetime
    )
  `)
  if err != nil {
    log.Print(err)
  }

  _,err = db.Exec(`
    create table scheduled_posts(
      id integer primary key autoincrement,
      img_id int, 
      text varchar, 
      site varchar,
      date datetime,
      foreign key (img_id)
        references images(id)
          on delete cascade
      )
  `)
  if err != nil {
    log.Print(err)
  }

  return true
}

func (conn Connection) AddImage(name string, path string, date time.Time) {
  fn := func (db *sql.DB) {
    db.Exec(`
      insert into images(path,name,date)
      values(?, ?, ?)
    `,
    path,
    name,
    date,
    )
  }

  conn.Do(fn)
}

func (conn Connection) GetAllImages() (images []Image) {
  fn := func (db *sql.DB) {
    rows,err := db.Query(`
      select id,name,path,date from images 
    `)
    if err != nil {
      log.Println(err)
      return
    }

    for rows.Next() {
      var image Image
      if err := rows.Scan(&image.Id, &image.Name, &image.Path, &image.Date) ; err != nil {
        log.Println(err)
        continue
      }
      images = append(images, image)
    }
  }

  conn.Do(fn)
  return
}

func (conn Connection) GetImage(id int) (image Image) {
  fn := func (db *sql.DB) {
    image.Id = id
    image.Name = ""
    image.Path = ""
    image.Date = time.Now()

    row := db.QueryRow(`
      select name,path,date from images where id=?
    `,
    id,
    )
    row.Scan(&image.Name,&image.Path,&image.Date)
  }

  conn.Do(fn)
  return
}

func (conn Connection) RemoveImage(id int) {
  fn := func (db *sql.DB) {
    db.Exec(`
      delete from images
      where id=?
    `, id)
  }

  conn.Do(fn)
}

func (conn Connection) SchedulePost(image Image, text string, site string, date time.Time) {
  fn := func (db *sql.DB) {
    db.Exec(`
      insert into scheduled_posts(img_id,text,site,date)
      values(?,?,?,?)
    `,
    image.Id,
    text,
    site,
    date,
    )
  }

  conn.Do(fn)
}

func (conn Connection) RemoveScheduledPost(id int) {
  fn := func (db *sql.DB) {
    db.Exec(`
      delete from scheduled_posts
      where id=?
    `, id)
  }

  conn.Do(fn)
}

func (conn Connection) GetEarliestPost() (post ScheduledPost) {
  fn := func (db *sql.DB) {
    row := db.QueryRow(`
      select (id,img_id,text,site,date) from schedule_posts
      order by date
      limit 1
    `)

    row.Scan(&post.Id,&post.ImageId,&post.Text,&post.Site,&post.Date)
  }

  conn.Do(fn)
  return
}

func (conn Connection) GetAllScheduledPosts() (posts []ScheduledPost) {
  fn := func(db *sql.DB) {
    rows,_ := db.Query(`
      select id,img_id,text,site,date from scheduled_posts
    `)
    for rows.Next() {
      var post ScheduledPost
      if err := rows.Scan(&post.Id,&post.ImageId,&post.Text,&post.Site,&post.Date) ; err != nil {
        log.Println(err)
        continue
      }
      posts = append(posts, post)
    }
  }

  conn.Do(fn)
  return
}
