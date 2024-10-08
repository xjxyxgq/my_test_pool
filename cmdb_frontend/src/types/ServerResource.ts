export interface ServerResource {
  id: number;
  cluster_id: number;
  ip: string;
  clusterName: string;
  instance_id: string;
  group_name: string;
  cluster_name: string;
  used_memory: number;
  total_memory: number;
  used_disk: number;
  total_disk: number;
  cpu_load: number;
  cpuUsage: number;
  memoryUsage: number;
  diskUsage: number;
  cpuThreshold: number;
  memoryThreshold: number;
  diskThreshold: number;
  pool_id: number;
  port: number;
  instance_role: string;
  cpu_cores: number;
  date_time: string;
  department_name: string;
}