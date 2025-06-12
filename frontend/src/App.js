import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/common/Layout';
import AssetList from './components/assets/AssetList';
import AssetDetail from './components/assets/AssetDetail';
import AssetForm from './components/assets/AssetForm';
import ServiceList from './components/services/ServiceList';
import ServiceForm from './components/services/ServiceForm';
import LogList from './components/logs/LogList';
import './App.css';

function App() {
  return (
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<AssetList />} />
            <Route path="/assets" element={<AssetList />} />
            <Route path="/assets/new" element={<AssetForm />} />
            <Route path="/assets/:ip" element={<AssetDetail />} />
            <Route path="/assets/:ip/edit" element={<AssetForm />} />
            <Route path="/assets/:ip/services" element={<ServiceList />} />
            <Route path="/assets/:ip/services/new" element={<ServiceForm />} />
            <Route path="/assets/:ip/services/:serviceId/edit" element={<ServiceForm />} />
            <Route path="/logs" element={<LogList />} />
          </Routes>
        </Layout>
      </Router>
  );
}

export default App;