Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Id_cat,
Title_post,
Content,
Likes INTEGER,
Posted_user,
Last_Posted,
Nb_Reply

--Field on 11-12

Id_cat,Title_post,Content,Likes,Posted_user,Last_Posted,Nb_Reply

--test_post on 15-16

'1','recette','recette_content','0','Naywvi','0','0'



--FOREIGN KEY (Id_cat) REFERENCES categorie(Id),
--FOREIGN KEY (Posted_user) REFERENCES user(Id)