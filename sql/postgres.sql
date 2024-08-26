
-- 用户
create table sys_user(
    id bigserial primary key,
    role int not null,
    department int,
    username varchar unique,
    name varchar,
    phone varchar unique,
    pwd varchar,
    gender smallint  DEFAULT 0 not null, -- 1-nan,2-nv
    image varchar,   -- 头像
    address varchar,
    status smallint DEFAULT 0 not null,  -- 冻结 1
    deleted boolean default false,
    extend jsonb default '{}'::jsonb, -- 权限剔除privilegeExclude:[]; 编号：no； 岗位 post
    createDt timestamp default (now() at time zone 'PRC')
);

-- 角色
create table sys_role(
    id bigserial primary key,
    department int default 0,
    name varchar not null,
    description varchar,
    privileges text[],
    createDt timestamp default (now() at time zone 'PRC'),
    deleted boolean default false,
    extend jsonb default '{}'::jsonb  -- immutable:不可删除
);

-- 权限常量
create table sys_privilege_constant(
    id varchar primary key,
    name varchar not null,
    type varchar, -- 分类
    sort int default 0
);

-- 部门
create table sys_department(
   id bigserial primary key,
   no varchar,	-- 编号
   name varchar,
   descr varchar,	-- 描述
   parent int,
   extend jsonb default '{}'::jsonb, -- 简称 refer；颜色-color
   createDt timestamp default (now() at time zone 'PRC'),
   deleted boolean default false
);

create table more_setting(
    id int primary key,
    data jsonb default '{}'::jsonb
);
insert into more_setting(id,data) values(1,'{}');
