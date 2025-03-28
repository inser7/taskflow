import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import LoginPage from './components/LoginPage';
import TaskList from './components/TaskList';

const App: React.FC = () => {
  // Простая функция проверки аутентификации
  const isAuthenticated = () => {
    return localStorage.getItem('token') !== null;
  };

  // Защищенный роут, который редиректит на логин, если нет токена
  const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return isAuthenticated() ? <>{children}</> : <Navigate to="/login" replace />;
  };

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route 
            path="/tasks" 
            element={
              <PrivateRoute>
                <TaskList />
              </PrivateRoute>
            } 
          />
          {/* Редирект по умолчанию */}
          <Route 
            path="/" 
            element={<Navigate to={isAuthenticated() ? "/tasks" : "/login"} replace />} 
          />
        </Routes>
      </div>
    </Router>
  );
};

export default App;