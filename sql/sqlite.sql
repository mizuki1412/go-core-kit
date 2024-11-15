
-- 用户
create table sys_user(
    id integer primary key,
    role int not null,
    department int,
    username varchar unique,
    name varchar,
    phone varchar unique,
    pwd varchar,
    gender int DEFAULT 0 not null, -- 1-nan,2-nv
    image varchar,   -- 头像
    address varchar,
    status smallint DEFAULT 0 not null,  -- 冻结 1
    deleted boolean default false,
    extend text, -- 权限剔除privilegeExclude:[]; 编号：no； 岗位 post
    createdt integer
);

-- 角色
create table sys_role(
    id integer primary key,
    department int default 0,
    name varchar not null,
    description varchar,
    privileges text,
    createdt integer,
    deleted boolean default false,
    extend text  -- immutable:不可删除
);

-- 权限常量
create table sys_privilege_constant(
    id text primary key,
    name varchar not null,
    type varchar, -- 分类
    sort int default 0
);

-- 部门
create table sys_department(
   id integer primary key,
   no varchar,	-- 编号
   name varchar,
   descr varchar,	-- 描述
   parent int,
   extend text, -- 简称 refer；颜色-color
   createdt integer,
   deleted boolean default false
);

create table more_setting(
    id integer primary key,
    data text
);
insert into more_setting(id,data) values(1,'{}');
