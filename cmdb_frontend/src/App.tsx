import React, { useState, useEffect } from 'react';
import { Layout, Tabs } from 'antd';
import HostList from './components/HostList';
import DatabaseClusterAnalysis from './components/DatabaseClusterAnalysis';
import './App.css';

const { Header, Content } = Layout;
const { TabPane } = Tabs;

const App: React.FC = () => {
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, []);

  const formatTime = (date: Date) => {
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).replace(/\//g, '-');
  };

  return (
    <Layout className="layout">
      <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ color: 'white', margin: 0 }}>数据库主机资源池</h1>
        <span style={{ color: 'white' }}>{formatTime(currentTime)}</span>
      </Header>
      <Content style={{ padding: '0 50px' }}>
        <div className="site-layout-content">
          <Tabs defaultActiveKey="1">
            <TabPane tab="主机资源池" key="1">
              <HostList />
            </TabPane>
            <TabPane tab="主机资源用量分析" key="2">
              <DatabaseClusterAnalysis />
            </TabPane>
          </Tabs>
        </div>
      </Content>
    </Layout>
  );
}

export default App;