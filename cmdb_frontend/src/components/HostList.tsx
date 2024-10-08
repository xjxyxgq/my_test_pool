import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { Table, Input, Select, Button, Space, Modal } from 'antd';
import HostDetail from './HostDetail';

// 更新 Host 接口
interface Host {
  id: number;
  host_name: string;
  host_ip: string;
  host_type: string;
  h3c_id: string;
  h3c_status: string;
  disk_size: number;
  ram: number;
  vcpus: number;
  if_h3c_sync: string;
  h3c_img_id: string;
  h3c_hm_name: string;
  is_delete: string;
  leaf_number: string;
  rack_number: string;
  rack_height: number;
  rack_start_number: number;
  from_factor: number;
  serial_number: string;
  is_deleted: boolean;
  is_static: boolean;
  create_time: string;
  update_time: string;
  host_applications: Application[];
}

// 更新 Application 接口
interface Application {
  id: number;
  pool_id: number;
  server_type: string;
  server_version: string;
  server_subtitle: string;
  cluster_name: string;
  server_protocol: string;
  server_addr: string;
  server_port: number;
  server_role: string;
  server_status: string;
  department_name: string;
  create_time: string;
  update_time: string;
}

const { Option } = Select;

const HostList: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [filteredHosts, setFilteredHosts] = useState<Host[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [ipFilter, setIpFilter] = useState('');
  const [datacenterFilter, setDatacenterFilter] = useState<string[]>([]);
  const [appTypeFilter, setAppTypeFilter] = useState<string[]>([]);
  const [departmentFilter, setDepartmentFilter] = useState<string[]>([]);
  const [selectedHost, setSelectedHost] = useState<Host | null>(null);
  const [isModalVisible, setIsModalVisible] = useState(false);

  const getIDCNameFromIP = (ip: string): string => {
    const parts = ip.split('.');
    if (parts.length >= 2) {
      switch (parts[1]) {
        case '1':
          return 'P1';
        case '2':
          return 'P2';
        case '3':
          return 'P3';
        case '4':
          return 'P4';
        case '5':
          return 'P5';
        case '6':
          return 'P6';
        default:
          return 'Unknown-IDC';
      }
    }
    return 'Unknown-IDC';
  };

  useEffect(() => {
    setLoading(true);
    axios.get('/api/cmdb/v1/get_hosts_pool_detail')
      .then(response => {
        setHosts(response.data);
        setFilteredHosts(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching hosts:', error);
        setError('Failed to fetch hosts data');
        setLoading(false);
        setHosts([]);
        setFilteredHosts([]);
      });
  }, []);

  const applyFilters = useCallback(() => {
    let filtered = hosts;

    if (ipFilter) {
      filtered = filtered.filter(host => host.host_ip.includes(ipFilter));
    }

    if (datacenterFilter.length > 0) {
      filtered = filtered.filter(host => {
        const idcName = getIDCNameFromIP(host.host_ip);
        return datacenterFilter.includes(idcName);
      });
    }

    if (appTypeFilter.length > 0) {
      filtered = filtered.filter(host => 
        host.host_applications.some(app => appTypeFilter.includes(app.server_type))
      );
    }

    if (departmentFilter.length > 0) {
      filtered = filtered.filter(host => 
        host.host_applications.some(app => departmentFilter.includes(app.department_name))
      );
    }

    setFilteredHosts(filtered);
  }, [hosts, ipFilter, datacenterFilter, appTypeFilter, departmentFilter]);

  useEffect(() => {
    applyFilters();
  }, [ipFilter, datacenterFilter, appTypeFilter, departmentFilter, hosts, applyFilters]);

  const showHostDetail = (host: Host) => {
    setSelectedHost(host);
    setIsModalVisible(true);
  };

  const columns = [
    {
      title: '主机名',
      dataIndex: 'host_name',
      key: 'host_name',
      sorter: (a: Host, b: Host) => a.host_name.localeCompare(b.host_name),
      render: (text: string, record: Host) => (
        <Button type="link" onClick={() => showHostDetail(record)}>
          {text}
        </Button>
      ),
    },
    {
      title: 'IP地址',
      dataIndex: 'host_ip',
      key: 'host_ip',
      sorter: (a: Host, b: Host) => a.host_ip.localeCompare(b.host_ip),
    },
    {
      title: 'CPU核数',
      dataIndex: 'vcpus',
      key: 'vcpus',
      sorter: (a: Host, b: Host) => a.vcpus - b.vcpus,
    },
    {
      title: '内存大小(GB)',
      dataIndex: 'ram',
      key: 'ram',
      sorter: (a: Host, b: Host) => a.ram - b.ram,
    },
    {
      title: '硬盘空间(GB)',
      dataIndex: 'disk_size',
      key: 'disk_size',
      sorter: (a: Host, b: Host) => a.disk_size - b.disk_size,
    },
    {
      title: '主机类型',
      dataIndex: 'host_type',
      key: 'host_type',
      sorter: (a: Host, b: Host) => a.host_type.localeCompare(b.host_type),
      render: (text: string) => text === '0' ? '云主机' : '裸金属',
    },
  ];

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input
          placeholder="IP地址"
          value={ipFilter}
          onChange={(e) => setIpFilter(e.target.value)}
          style={{ width: 200 }}
        />
        <Select
          mode="multiple"
          placeholder="选择机房"
          value={datacenterFilter}
          onChange={(value) => setDatacenterFilter(value)}
          style={{ width: 200 }}
        >
          <Option value="P1">P1</Option>
          <Option value="P2">P2</Option>
          <Option value="P3">P3</Option>
          <Option value="P4">P4</Option>
          <Option value="P5">P5</Option>
          <Option value="P6">P6</Option>
        </Select>
        <Select
          mode="multiple"
          placeholder="应用类型"
          value={appTypeFilter}
          onChange={(value) => setAppTypeFilter(value)}
          style={{ width: 200 }}
        >
          {Array.from(new Set(hosts.flatMap(host => 
            host.host_applications ? host.host_applications.map(app => app.server_type) : []
          ))).map(type => (
            <Option key={type} value={type}>{type}</Option>
          ))}
        </Select>
        <Select
          mode="multiple"
          placeholder="所属部门"
          value={departmentFilter}
          onChange={(value) => setDepartmentFilter(value)}
          style={{ width: 200 }}
        >
          {Array.from(new Set(hosts.flatMap(host => 
            host.host_applications ? host.host_applications.map(app => app.department_name) : []
          ))).map(dept => (
            <Option key={dept} value={dept}>{dept}</Option>
          ))}
        </Select>
        <Button onClick={() => {
          setIpFilter('');
          setDatacenterFilter([]);
          setAppTypeFilter([]);
          setDepartmentFilter([]);
        }}>重置</Button>
      </Space>
      <Table
        columns={columns}
        dataSource={filteredHosts.map(host => ({ ...host, key: host.id }))}
        rowKey="id"
        loading={loading}
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['5', '10', '20', '50'],
          defaultPageSize: 5,
        }}
        onChange={(pagination, filters, sorter) => {
          // 如果需要，可以在这里处理排序变化
          console.log('sorter', sorter);
        }}
      />
      <Modal
        title="主机详情"
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
        width={800}
      >
        {selectedHost && <HostDetail host={selectedHost} />}
      </Modal>
    </div>
  );
};

export default HostList;