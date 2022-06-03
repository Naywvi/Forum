Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
Name,
Create_post,
Del_profile,
Signal_post,
Moove_post,
Comment_post,
Del_post,
Del_post_user,
Admin

--Field in 14-13

Name, Create_post, Del_profile, Signal_post, Moove_post, Comment_post, Del_post, Del_post_user, Admin

--Rank with perm | Admin -> 17-18 | Modo 19-20 | User 21-22 | Default 23-24

'Admin', '1', '1', '1', '1', '1', '1', '1', '1'

'Modo', '1', '1', '1', '1', '1', '1', '1', '0'

'User', '1', '0', '1', '0', '1', '1', '0', '0'

'Default', '0', '0', '1', '0', '0', '0', '0', '0'