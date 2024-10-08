import React from 'react';
// eslint-disable-next-line
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Row, Col, Card, Table } from 'antd';

// 抑制 YAxis 和 XAxis 的 defaultProps 警告
YAxis.defaultProps = {
  ...YAxis.defaultProps,
  allowDataOverflow: true,
};

XAxis.defaultProps = {
  ...XAxis.defaultProps,
  allowDataOverflow: true,
};

export interface ClusterResource {
  clusterName: string;
  groupName: string;
  memory: number;
  memoryTotal: number;
  memoryUsed: number;
  disk: number;
  diskTotal: number;
  diskUsed: number;
  cpu: number;
  cpuTotal: number;
  cpuUsed: number;
  count: number;
  maxMemory: number; // 新增字段
  maxDisk: number;   // 新增字段
  maxCPU: number;    // 新增字段
}

interface ClusterResourceUsageProps {
  data: ClusterResource[];
}

const ClusterResourceUsage: React.FC<ClusterResourceUsageProps> = ({ data }) => {
  const renderClusterChart = (cluster: ClusterResource, index: number) => {
    const chartData = [
      { 
        name: '内存', 
        平均使用率: cluster.memory.toFixed(2), 
        最大使用率: cluster.maxMemory.toFixed(2), 
        total: cluster.memoryTotal, 
        used: cluster.memoryUsed 
      },
      { 
        name: '磁盘', 
        平均使用率: cluster.disk.toFixed(2), 
        最大使用率: cluster.maxDisk.toFixed(2), 
        total: cluster.diskTotal, 
        used: cluster.diskUsed 
      },
      { 
        name: 'CPU', 
        平均使用率: cluster.cpu.toFixed(2), 
        最大使用率: cluster.maxCPU.toFixed(2), 
        total: cluster.cpuTotal, 
        used: cluster.cpuUsed 
      }
    ];

    const CustomTooltip = ({ active, payload, label }: any) => {
      if (active && payload && payload.length) {
        const data = payload[0].payload;
        return (
          <div style={{ backgroundColor: 'white', padding: '5px', border: '1px solid #ccc' }}>
            <p>{`${label} : ${parseFloat(data.平均使用率).toFixed(2)}% (平均)`}</p>
            <p>{`${label} : ${parseFloat(data.最大使用率).toFixed(2)}% (最大)`}</p>
            {label === 'CPU' ? (
              <p>{`使用量: ${data.used.toFixed(2)}%`}</p>
            ) : (
              <p>{`使用量: ${data.used.toFixed(2)} GB / ${data.total.toFixed(2)} GB`}</p>
            )}
          </div>
        );
      }
      return null;
    };

    return (
      <Col key={`cluster-${index}`} style={{ width: 300, marginBottom: '20px' }}>
        <Card 
          title={`${cluster.groupName}-${cluster.clusterName}`} 
          style={{ width: 300, height: 400, overflow: 'hidden' }}
          styles={{
            body: { height: 'calc(100% - 57px)', padding: '12px' }
          }}
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={chartData}
              margin={{
                top: 5,
                right: 10, // 保留右侧10px空白
                left: 10,  // 保留左侧10px空白
                bottom: 5,
              }}
              barCategoryGap={30} // 设置柱状图之间的间隔
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey="name" 
                scale="point" 
                padding={{ left: 40, right: 40 }}
                tick={{ fontSize: 12 }}
                allowDataOverflow={true}
              />
              <YAxis 
                domain={[0, 100]} 
                allowDataOverflow={true}
                tick={{ fontSize: 12 }}
              />
              <Tooltip content={<CustomTooltip />} />
              <Legend wrapperStyle={{ fontSize: 12 }} />
              <Bar dataKey="平均使用率" fill="#8884d8" barSize={15} />
              <Bar dataKey="最大使用率" fill="#82ca9d" barSize={15} />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      </Col>
    );
  };

  const columns = [
    {
      title: '集群名称',
      dataIndex: 'clusterName',
      key: 'clusterName',
      sorter: (a: ClusterResource, b: ClusterResource) => a.clusterName.localeCompare(b.clusterName),
    },
    {
      title: '组名',
      dataIndex: 'groupName',
      key: 'groupName',
      sorter: (a: ClusterResource, b: ClusterResource) => a.groupName.localeCompare(b.groupName),
    },
    {
      title: '内存使用率 (平均)',
      dataIndex: 'memory',
      key: 'memory',
      sorter: (a: ClusterResource, b: ClusterResource) => a.memory - b.memory,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
    {
      title: '内存使用率 (最大)',
      dataIndex: 'maxMemory',
      key: 'maxMemory',
      sorter: (a: ClusterResource, b: ClusterResource) => a.maxMemory - b.maxMemory,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
    {
      title: '磁盘使用率 (平均)',
      dataIndex: 'disk',
      key: 'disk',
      sorter: (a: ClusterResource, b: ClusterResource) => a.disk - b.disk,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
    {
      title: '磁盘使用率 (最大)',
      dataIndex: 'maxDisk',
      key: 'maxDisk',
      sorter: (a: ClusterResource, b: ClusterResource) => a.maxDisk - b.maxDisk,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
    {
      title: 'CPU使用率 (平均)',
      dataIndex: 'cpu',
      key: 'cpu',
      sorter: (a: ClusterResource, b: ClusterResource) => a.cpu - b.cpu,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
    {
      title: 'CPU使用率 (最大)',
      dataIndex: 'maxCPU',
      key: 'maxCPU',
      sorter: (a: ClusterResource, b: ClusterResource) => a.maxCPU - b.maxCPU,
      render: (value: number) => `${value.toFixed(2)}%`,
    },
  ];

  return (
    <Row gutter={16}>
      <Col span={24}>
        <Table columns={columns} dataSource={data} rowKey={(record) => `${record.groupName}-${record.clusterName}`} pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['5', '10', '20', '50'],
          defaultPageSize: 5,
        }} />
      </Col>
      {data.map((cluster, index) => renderClusterChart(cluster, index))}
    </Row>
  );
};

export default ClusterResourceUsage;
