import React from 'react';
import { Table, TablePaginationConfig } from 'antd';
import { ServerResource } from '../types/ServerResource';

interface DiskFullPredictionProps {
  data: ServerResource[];
  pagination?: TablePaginationConfig; // 添加 pagination 属性
}

const DiskFullPrediction: React.FC<DiskFullPredictionProps> = ({ data, pagination }) => {
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
      title: 'Current Disk Usage',
      dataIndex: 'diskUsage',
      key: 'diskUsage',
      sorter: (a: ServerResource, b: ServerResource) => (a.used_disk / a.total_disk) - (b.used_disk / b.total_disk),
      render: (text: number, record: ServerResource) => `${((record.used_disk / record.total_disk) * 100).toFixed(2)}%`,
    },
    {
      title: 'Predicted Full Date',
      dataIndex: 'predictedFullDate',
      key: 'predictedFullDate',
      sorter: (a: ServerResource, b: ServerResource) => {
        const getFullDate = (record: ServerResource) => {
          const usageRate = (record.used_disk / record.total_disk) / 30;
          const daysUntilFull = (1 - (record.used_disk / record.total_disk)) / usageRate;
          return new Date(new Date().getTime() + daysUntilFull * 24 * 60 * 60 * 1000);
        };
        return getFullDate(a).getTime() - getFullDate(b).getTime();
      },
      render: (text: string, record: ServerResource) => {
        const usageRate = (record.used_disk / record.total_disk) / 30;
        const daysUntilFull = (1 - (record.used_disk / record.total_disk)) / usageRate;
        const fullDate = new Date(new Date().getTime() + daysUntilFull * 24 * 60 * 60 * 1000);
        return fullDate.toLocaleDateString();
      },
    },
  ];

  return (
    <Table columns={columns} dataSource={data} rowKey={(record) => `${record.instance_id}-${record.ip}`} pagination={{
      showSizeChanger: true,
      showQuickJumper: true,
      pageSizeOptions: ['5', '10', '20', '50'],
      defaultPageSize: 5,
    }} />
  );
}

export default DiskFullPrediction;