import React, { useEffect, useState } from 'react';
import { Alert } from 'antd';
import { ServerResource } from '../types/ServerResource';

interface AlertsProps {
  data: ServerResource[];
  lowThreshold: number;
  highThreshold: number;
  triggerUpdate: number;
  selectedGroups: string[];
}

const Alerts: React.FC<AlertsProps> = ({ data, lowThreshold, highThreshold, triggerUpdate, selectedGroups }) => {
  const [alerts, setAlerts] = useState<JSX.Element[]>([]);

  useEffect(() => {
    const filteredData = selectedGroups.length > 0
      ? data.filter(resource => selectedGroups.includes(resource.group_name))
      : data;

    const newAlerts = filteredData.flatMap((item, index) => {
      const memoryUsage = (item.used_memory / item.total_memory) * 100;
      const diskUsage = (item.used_disk / item.total_disk) * 100;
      const cpuUsage = item.cpu_load;

      const alertPrefix = `${item.ip} (${item.group_name} ${item.cluster_name})`;
      const alerts: JSX.Element[] = [];

      // 使用 index 作为 key 的一部分，确保唯一性
      const uniqueKey = `${index}-${item.ip}`;

      // Memory alerts
      if (memoryUsage > highThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-memory-high`}
            message={`${alertPrefix} | 内存: ${memoryUsage.toFixed(2)}% (${item.used_memory.toFixed(2)}GB/${item.total_memory.toFixed(2)}GB) | 警告：高于${highThreshold}%阈值`}
            type="error"
            showIcon
            banner
          />
        );
      } else if (memoryUsage < lowThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-memory-low`}
            message={`${alertPrefix} | 内存: ${memoryUsage.toFixed(2)}% (${item.used_memory.toFixed(2)}GB/${item.total_memory.toFixed(2)}GB) | 提示：低于${lowThreshold}%阈值`}
            type="warning"
            showIcon
            banner
          />
        );
      }

      // Disk alerts
      if (diskUsage > highThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-disk-high`}
            message={`${alertPrefix} | 磁盘: ${diskUsage.toFixed(2)}% (${item.used_disk.toFixed(2)}GB/${item.total_disk.toFixed(2)}GB) | 警告：高于${highThreshold}%阈值`}
            type="error"
            showIcon
            banner
          />
        );
      } else if (diskUsage < lowThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-disk-low`}
            message={`${alertPrefix} | 磁盘: ${diskUsage.toFixed(2)}% (${item.used_disk.toFixed(2)}GB/${item.total_disk.toFixed(2)}GB) | 提示：低于${lowThreshold}%阈值`}
            type="warning"
            showIcon
            banner
          />
        );
      }

      // CPU alerts
      if (cpuUsage > highThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-cpu-high`}
            message={`${alertPrefix} | CPU: ${cpuUsage.toFixed(2)}% | 警告：高于${highThreshold}%阈值`}
            type="error"
            showIcon
            banner
          />
        );
      } else if (cpuUsage < lowThreshold) {
        alerts.push(
          <Alert
            key={`${uniqueKey}-cpu-low`}
            message={`${alertPrefix} | CPU: ${cpuUsage.toFixed(2)}% | 提示：低于${lowThreshold}%阈值`}
            type="warning"
            showIcon
            banner
          />
        );
      }

      return alerts;
    });

    setAlerts(newAlerts);
  }, [data, lowThreshold, highThreshold, triggerUpdate, selectedGroups]);

  return <>{alerts}</>;
};

export default Alerts;