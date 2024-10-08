import React from 'react';
import { Descriptions } from 'antd';

interface HostPool {
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
  host_applications: HostApplication[];
}

interface HostApplication {
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
  department_name: string;  // 新增字段
  create_time: string;
  update_time: string;
}

interface HostDetailProps {
  host: HostPool;
}

const HostDetail: React.FC<HostDetailProps> = ({ host }) => {
  const formatValue = (key: string, value: any) => {
    if (key === 'host_type') {
      return value === '0' ? 'Cloud Host' : 'Bare Metal';
    }
    if (key === 'is_deleted' || key === 'is_static') {
      return value ? 'Yes' : 'No';
    }
    if (key === 'create_time' || key === 'update_time') {
      return new Date(value).toLocaleString();
    }
    return value?.toString() || 'N/A';
  };

  return (
    <>
      <Descriptions title="Host Details" bordered column={3}>
        {Object.entries(host).map(([key, value]) => {
          if (key !== 'host_applications') {
            return (
              <Descriptions.Item key={`host-${key}`} label={key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}>
                {formatValue(key, value)}
              </Descriptions.Item>
            );
          }
          return null;
        })}
      </Descriptions>
      
      <h3 style={{ marginTop: '20px' }}>Host Applications</h3>
      {host.host_applications.map((app, index) => (
        <Descriptions key={`app-${app.id || index}`} title={`Application ${index + 1}`} bordered column={2} style={{ marginTop: '10px' }}>
          {Object.entries(app).map(([key, value]) => (
            <Descriptions.Item key={`app-${app.id || index}-${key}`} label={key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}>
              {formatValue(key, value)}
            </Descriptions.Item>
          ))}
        </Descriptions>
      ))}
    </>
  );
};

export default HostDetail;