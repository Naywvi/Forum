Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Id_cat INTEGER NOT NULL,
Title_post TEXT NOT NULL,
Content TEXT NOT NULL,
Likes INTEGER NOT NULL,
Posted_user INTEGER NOT NULL

--Field on 11-12

Id_cat,Title_post,Content,Likes,Posted_user

--test_post on 15-16

'1','recette','recette_content','0','Naywvi'



--FOREIGN KEY (Id_cat) REFERENCES categorie(Id),
--FOREIGN KEY (Posted_user) REFERENCES user(Id)