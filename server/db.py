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
        # delete any reference to the image in the scheduled post table
        # self.conn.execute('''
        #     delete from scheduled_posts 
        #     where img_id=:img_id
        #     ''', ({'img_id': img_id}))
        self.conn.commit()
        pass
    
    def schedule_post(self, img_id : int, text : str, date : datetime):
        self.conn.execute('''
            insert into scheduled_posts(img_id,text,date)
            values(:img_id,:text,:date)
        ''',(
            {'img_id': img_id, 'text': text, 'date': date.isoformat()}
        ))
        self.conn.commit()
        pass

    def get_earliest_post(self):
        for row in self.conn.execute('''
            select *
            from scheduled_posts
            order by date
            limit 1
        '''):
            return row
