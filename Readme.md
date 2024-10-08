## 一、需求
开发一个web应用程序，提供一个页面，展示资源池表里的数据库服务器信息，包括主机名、ip、cpu等该主机的主要软硬件信息。并且可以通过弹出窗口，展示更详细的服务器信息，包括所有具体的硬件信息，以及一个列表，展示所有属于该主机的软件信息。

## 二、开发架构和语言
1. 采用前后端分离的方式开发
2. 前端开发框架为 React
3. 后端开发语言为 Golang，通过 RESTful API 的方式为前端提供数据

## 三、后端数据字典
1. 资源池表，这里存储了所有的数据库服务器硬件信息
    ```sql
    create table `hosts_pool`(
        `id` bigint unsigned not null auto_increment comment '主键',
        `host_name` varchar(50) not null comment '主机名',
        `host_ip` varchar(50) not null comment '主机ip',
        `host_type` varchar(10) default null comment '0 云主机，1 裸金属',
        `h3c_id` varchar(50) default null comment 'h3c id',
        `h3c_status` varchar(20) default null comment 'h3c_status',
        `disk_size` int unsigned default null comment '硬盘空间',
        `ram` int unsigned default null comment '内存大小',
        `vcpus` int unsigned default null comment 'cpu核数',
        `if_h3c_sync` varchar(10) default null comment 'unknown',
        `h3c_img_id` varchar(50) default null comment 'img id',
        `h3c_hm_name` varchar(1000) default null comment 'hm name',
        `is_delete` varchar(10) default null comment '是否已删除',
        `leaf_number` varchar(50) default null comment '交换机编号',
        `rack_number` varchar(10) default null comment '机架编号',
        `rack_height` int unsigned default null comment '机架高度',
        `rack_start_number` int unsigned default null comment '机架起始编号',
        `from_factor` int unsigned default null comment '主机高度',
        `serial_number` varchar(50) default null comment '序列号',
        `is_deleted` tinyint not null default '0' comment '软删除标记',
        `is_static` tinyint not null default '0' comment '是否静态数据',
        `create_time` datetime not null default current_timestamp comment '写入时间',
        `update_time` datetime not null default current_timestamp on update current_timestamp comment '更新时间',
        primary key (`id`),
        unique key (`host_ip`)
    )
    ```

2. 主机的应用程序记录表，这里存储了所有主机上部署的应用程序
    ```sql
    create table `hosts_applications`(
        `id` bigint unsigned not null auto_increment comment '主键',
        `pool_id` bigint unsigned default not null comment '关联hosts_pool表的主键',
        `server_type` varchar(30) default null comment 'server type',
        `server_version` varchar(30) default null comment 'server version',
        `server_subtitle` varchar(30) default null comment '子属性',
        `cluster_group_name` varchar(64) default null comment '所属集群组名称',
        `cluster_name` varchar(64) default null comment '所属集群名称',
        `server_protocol` varchar(64) default null comment '服务的访问协议，如mysql、http等',
        `server_addr` varchar(100) default null comment '服务的访问地址',
        `department_name` varchar(100) default null comment '所属部门',
        `create_time` datetime not null default current_timestamp comment '写入时间',
        `update_time` datetime not null default current_timestamp on update current_timestamp comment '更新时间',
        primary key (`id`),
        unique key (`pool_id`, `service_addr`, `service_protocol`)
    )
    ```

## 四、后端 RESTful API 接口
1. `/api/cmdb/v1/get_hosts_pool_detail/`
    - 返回所有 `hosts_pool` 表中的主机，以及记录在 `hosts_applications` 表中所有的应用信息，一个主机可能部署多个应用。
2. `/api/cmdb/v1/collect_applications/`
    - 加载上述两个数据表的数据，模拟数据：
        1. `hosts_pool` 表模拟30条，表示30台服务器
        2. `hosts_applications` 表模拟50条，将这些应用程序分配给 `hosts_pool` 表模拟的服务器，每台分配0个或多个应用程序
3. `/api/cmdb/v1/get_host_detail/:id`
    - 根据主机ID获取主机的详细信息，包括硬件信息和应用信息。
4. `/api/cmdb/v1/get_application_detail/:id`
    - 根据应用ID获取应用的详细信息。
5. `/api/cmdb/v1/send_email`
    - 发送包含服务器资源使用情况报告的邮件。

## 五、前端页面
目前只需要一个主页面，主页面需要有这几个部分：
1. 页面的标题"服务器资源池"
2. 数据筛选部分，包括几个下拉框、文本框用于填充筛选条件，还包括一个重置按钮，用于对页面显示的数据进行筛选，可用的筛选项目包括：
    1. 文本框，输入ip地址进行模糊搜索。筛选在输入时实时触发。
    2. 下拉框，选择服务器所属的机房，包括 P1 到 P6 几个机房，服务器的前两位IP地址表示了它所属的机房，192.1 表示P1,192.2 表示P2，192.3 表示P3，192.4 表示P4，192.5 表示P5，192.6 表示P6。根据这个规则对数据进行筛选。筛选在选项变化时触发。支持多选。
    3. 下拉框，列出所有服务器上的应用类型并通过它过滤数据，服务器的应用类型来自于 `hosts_applications` 表中的 `server_type`。筛选在选项变化时触发。支持多选。
    4. 下拉框，列出所有服务器上的应用所属的部门并通过它过滤数据， 服务器的应用所属部门来自于 `hosts_applications` 表中的 `department_name`。筛选在选项变化时触发。支持多选。
3. 数据显示部分，通过一个分页列表展示数据，列表支持自定义单页大小，默认单页显示10条数据，所有数据列均支持排序。
4. 可以通过点击主机名，在弹出窗口中显示更详细的服务器信息，包括所有具体的硬件信息，以及一个列表，展示所有属于该主机的软件信息。

## 六、新增功能
1. 下拉框筛选支持多选：用户可以在机房、应用类型和部门的下拉框中选择多个选项进行筛选。
2. 主机详情页使用弹出窗口显示：点击主机名时，会弹出一个模态窗口显示该主机的详细信息，而不是跳转到新页面。

## 七、前端特性
1. 使用 Ant Design 组件库进行 UI 开发。
2. 使用 `axios` 进行 HTTP 请求。
3. 使用 `dayjs` 进行日期处理和本地化。
4. 支持页面截图并通过邮件发送功能。

## 八、后端特性
1. 使用 Gin Web Framework 进行 API 开发。
2. 使用 GORM 进行数据库操作。
3. 支持 CORS。
4. 支持通过 SMTP 发送邮件。

## 九、用法
### 前端
1. 安装依赖：
    ```bash
    npm install
    ```
2. 启动开发服务器：
    ```bash
    npm start
    ```

### 后端
1. 安装依赖：
    ```bash
    go mod tidy
    ```
2. 启动后端服务：
    ```bash
    go run main.go
    ```

## 十、贡献
欢迎贡献代码、提出建议或问题、修复 Bug 以及参与讨论对新功能的想法。请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解更多信息。

## 十一、许可证
本项目遵循 BSD 3-Clause 开源许可协议，访问 [https://opensource.org/licenses/BSD-3-Clause](https://opensource.org/licenses/BSD-3-Clause) 查看许可协议文件。

