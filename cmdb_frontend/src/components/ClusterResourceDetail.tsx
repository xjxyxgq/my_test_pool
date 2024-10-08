import React from 'react';
import { Table, TablePaginationConfig } from 'antd';
import { ServerResource } from '../types/ServerResource';

interface ClusterResourceDetailProps {
  data: ServerResource[];
  pagination?: TablePaginationConfig; // 添加 pagination 属性
}

const ClusterResourceDetail: React.FC<ClusterResourceDetailProps> = ({ data, pagination }) => {
  const columns = [
    { 
      title: 'ID', 
      dataIndex: 'id', 
      key: 'id',
      sorter: (a: ServerResource, b: ServerResource) => a.id - b.id,
    },
    { 
      title: 'Cluster Name', 
      dataIndex: 'cluster_name', 
      key: 'cluster_name',
      sorter: (a: ServerResource, b: ServerResource) => a.cluster_name.localeCompare(b.cluster_name),
    },
    { 
      title: 'Group Name', 
      dataIndex: 'group_name', 
      key: 'group_name',
      sorter: (a: ServerResource, b: ServerResource) => a.group_name.localeCompare(b.group_name),
    },
    { 
      title: 'IP', 
      dataIndex: 'ip', 
      key: 'ip',
      sorter: (a: ServerResource, b: ServerResource) => a.ip.localeCompare(b.ip),
    },
    { 
      title: 'Port', 
      dataIndex: 'port', 
      key: 'port',
      sorter: (a: ServerResource, b: ServerResource) => a.port - b.port,
    },
    { 
      title: 'Instance Role', 
      dataIndex: 'instance_role', 
      key: 'instance_role',
      sorter: (a: ServerResource, b: ServerResource) => a.instance_role.localeCompare(b.instance_role),
    },
    { 
      title: 'Total Memory (GB)', 
      dataIndex: 'total_memory', 
      key: 'total_memory', 
      render: (value: number) => value.toFixed(2),
      sorter: (a: ServerResource, b: ServerResource) => a.total_memory - b.total_memory,
    },
    { 
      title: 'Used Memory (GB)', 
      dataIndex: 'used_memory', 
      key: 'used_memory', 
      render: (value: number) => value.toFixed(2),
      sorter: (a: ServerResource, b: ServerResource) => a.used_memory - b.used_memory,
    },
    { 
      title: 'Total Disk (GB)', 
      dataIndex: 'total_disk', 
      key: 'total_disk', 
      render: (value: number) => value.toFixed(2),
      sorter: (a: ServerResource, b: ServerResource) => a.total_disk - b.total_disk,
    },
    { 
      title: 'Used Disk (GB)', 
      dataIndex: 'used_disk', 
      key: 'used_disk', 
      render: (value: number) => value.toFixed(2),
      sorter: (a: ServerResource, b: ServerResource) => a.used_disk - b.used_disk,
    },
    { 
      title: 'CPU Cores', 
      dataIndex: 'cpu_cores', 
      key: 'cpu_cores',
      sorter: (a: ServerResource, b: ServerResource) => a.cpu_cores - b.cpu_cores,
    },
    { 
      title: 'CPU Load (%)', 
      dataIndex: 'cpu_load', 
      key: 'cpu_load', 
      render: (value: number) => value.toFixed(2),
      sorter: (a: ServerResource, b: ServerResource) => a.cpu_load - b.cpu_load,
    },
    { 
      title: 'Date Time', 
      dataIndex: 'date_time', 
      key: 'date_time',
      sorter: (a: ServerResource, b: ServerResource) => new Date(a.date_time).getTime() - new Date(b.date_time).getTime(),
    },
  ];

  return <Table columns={columns} dataSource={data} rowKey={(record) => `${record.id}-${record.ip}`} pagination={{
    showSizeChanger: true,
    showQuickJumper: true,
    pageSizeOptions: ['5', '10', '20', '50'],
    defaultPageSize: 5,
  }}/>;
};

export default ClusterResourceDetail;