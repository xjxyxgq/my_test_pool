import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
// import 'antd/dist/antd.css'; // 如果使用 Antd 4.x 版本
import 'antd/dist/reset.css'; // 如果使用 Antd 5.x 版本

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
