import React, { useMemo } from 'react';
import { Table, TablePaginationConfig } from 'antd';
import { ServerResource } from '../types/ServerResource';

interface ResourceAlertsProps {
  data: ServerResource[];
  lowThreshold: number;
  highThreshold: number;
  triggerUpdate: number;
  pagination?: TablePaginationConfig; // 添加 pagination 属性
}

const ResourceAlerts: React.FC<ResourceAlertsProps> = ({ data, lowThreshold, highThreshold, triggerUpdate, pagination }) => {
  const filteredData = useMemo(() => {
    return data.filter(item => {
      const cpuUsage = item.cpu_load;
      const memoryUsage = (item.used_memory / item.total_memory) * 100;
      const diskUsage = (item.used_disk / item.total_disk) * 100;
      
      return cpuUsage < lowThreshold || cpuUsage > highThreshold ||
             memoryUsage < lowThreshold || memoryUsage > highThreshold ||
             diskUsage < lowThreshold || diskUsage > highThreshold;
    });
  }, [data, lowThreshold, highThreshold, triggerUpdate]);

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      sorter: (a: ServerResource, b: ServerResource) => a.id - b.id,
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip',
      sorter: (a: ServerResource, b: ServerResource) => a.ip.localeCompare(b.ip),
    },
    {
      title: 'Cluster Name',
      dataIndex: 'cluster_name',
      key: 'cluster_name',
      sorter: (a: ServerResource, b: ServerResource) => a.cluster_name.localeCompare(b.cluster_name),
    },
    {
      title: 'CPU Usage',
      dataIndex: 'cpuUsage',
      key: 'cpuUsage',
      sorter: (a: ServerResource, b: ServerResource) => a.cpu_load - b.cpu_load,
      render: (text: number, record: ServerResource) => {
        const usage = record.cpu_load;
        return <span style={{ color: usage < lowThreshold ? 'green' : usage > highThreshold ? 'red' : 'inherit' }}>{usage.toFixed(2)}%</span>;
      },
    },
    {
      title: 'Memory Usage',
      dataIndex: 'memoryUsage',
      key: 'memoryUsage',
      sorter: (a: ServerResource, b: ServerResource) => (a.used_memory / a.total_memory) - (b.used_memory / b.total_memory),
      render: (text: number, record: ServerResource) => {
        const usage = (record.used_memory / record.total_memory) * 100;
        return <span style={{ color: usage < lowThreshold ? 'green' : usage > highThreshold ? 'red' : 'inherit' }}>{usage.toFixed(2)}%</span>;
      },
    },
    {
      title: 'Disk Usage',
      dataIndex: 'diskUsage',
      key: 'diskUsage',
      sorter: (a: ServerResource, b: ServerResource) => (a.used_disk / a.total_disk) - (b.used_disk / b.total_disk),
      render: (text: number, record: ServerResource) => {
        const usage = (record.used_disk / record.total_disk) * 100;
        return <span style={{ color: usage < lowThreshold ? 'green' : usage > highThreshold ? 'red' : 'inherit' }}>{usage.toFixed(2)}%</span>;
      },
    },
  ];

  return (
    <Table columns={columns} dataSource={filteredData} rowKey={(record) => `${record.instance_id}-${record.ip}`} pagination={{
      showSizeChanger: true,
      showQuickJumper: true,
      pageSizeOptions: ['5', '10', '20', '50'],
      defaultPageSize: 5,
    }} />
  );
}

export default ResourceAlerts;