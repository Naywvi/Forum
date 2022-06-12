Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Id_post,
Date_comment,
User_posted,
Rank_User_Posted,
Title_comment,
Reply_user,
Reply_user_rank,
Reply_content,
Content_comment,
Likes

--Field on 11-12

Id_post,Date_comment,User_posted,Rank_User_Posted,Title_comment,Reply_user,Reply_user_rank,Reply_content,Content_comment,Likes

--test_post on 15-16

'1','01-0001-01','test','3','title','reply_to_test','reply_rank','<p>reply_content</p>','<p>content</p>','0'



--FOREIGN KEY (Id_post) REFERENCES categorie(Id),
--FOREIGN KEY (Posted_user) REFERENCES user(Id)