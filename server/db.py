import sqlite3
from datetime import datetime

class Connection:
    def __init__(self, path):
        self.conn = sqlite3.connect(path)
        
        self.conn.execute('pragma foreign_keys = on')
        # images db
        try:
            self.conn.execute('''
                create table images(
                    id integer primary key autoincrement, 
                    path varchar, 
                    name varchar, 
                    date datetime)''')
        except sqlite3.Error as e:
            print(e)

        # schedule
        try:
            self.conn.execute('''
                create table scheduled_posts(
                    id integer primary key autoincrement,
                    img_id int, 
                    text varchar, 
                    site varchar,
                    date datetime,
                    foreign key (img_id)
                        references images(id)
                            on delete cascade
                    )''')
        except sqlite3.Error as e:
            print(e)

    def add_img(self, path : str, name : str, date : datetime):
        data = (
            {'p': path, 'n': name, 'd': date.isoformat()}
        )
        self.conn.execute('''
            insert into images(path,name,date)
            values(:p, :n, :d)
            ''',
            data
        )
        self.conn.commit()
        pass
    
    def get_imgs(self) -> []:
        imgs = []
        for row in self.conn.execute('select * from images'):
            imgs.append(row)
        return imgs

    def remove_img(self, img_id : int):
        # delete image from image table
        self.conn.execute('''
            delete from images
            where id=:img_id
            ''', ({'img_id': img_id}))
        # the db restraints should delete any scheduled posts with this image
        self.conn.commit()
        pass
    
    def schedule_post(self, img_id : int, text : str, site : str, date : datetime):
        try:
            self.conn.execute('''
                insert into scheduled_posts(img_id,text,site,date)
                values(:img_id,:text,:site,:date)
            ''',(
                {'img_id': img_id, 'text': text, 'site': site, 'date': date.isoformat()}
            ))
            self.conn.commit()
        except sqlite3.Error as e:
            print(e)
        pass

    def remove_post(self, post_id):
        try:
            self.conn.execute('''
                delete from scheduled_post
                where id=:post_id
                ''',({'post_id': post_id}))
        except sqlite3.Error as e:
            print(e)
        pass


    def get_earliest_post(self):
        for row in self.conn.execute('''
            select *
            from scheduled_posts
            order by date
            limit 1
        '''):
            return row

    def get_scheduled_posts():
        posts = []
        for row in self.conn.execut('''
            select *
            from scheduled_posts
        '''):
            posts.append(row)
        return posts
