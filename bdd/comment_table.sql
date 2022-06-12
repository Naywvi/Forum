Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Id_post,
User,
Reply_to,
Content,
Likes

--Field on 11-12

Id_post,User,Reply_to,Content,Likes

--test_post on 15-16

'1','test','reply_to_me','none_comment','0'



--FOREIGN KEY (Id_post) REFERENCES categorie(Id),
--FOREIGN KEY (Posted_user) REFERENCES user(Id)