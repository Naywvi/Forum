Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Name,
Pswd,
Desc,
Email,
Profile_Picture,
Rank_id INTEGER NOT NULL,
FOREIGN KEY (Rank_id) REFERENCES rank(Id)

--Field on 11-12

Name,Pswd,Desc,Email,Profile_picture,Rank_id

--test_user on 16-16

'Naywvi','1230','none_dec','test@test.fr','none_picture',3