# YYUEsys
2020/04/01  部门数据        完成<br>
2020/04/13  员工数据        完成<br>
2020/04/13  班级数据        完成<br>
2020/04/15  服务项目        完成<br>
2020/04/25  合同内容        完成<br>



Mysql session表结构

    CREATE TABLE `session` (
        `session_key` char(64) NOT NULL,
        `session_data` blob,
        `session_expiry` int(11) unsigned NOT NULL,
        PRIMARY KEY (`session_key`)
    ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
