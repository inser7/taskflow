import React, { useEffect, useState } from 'react';

interface Task {
  id: number;
  name: string;
}

const TaskList: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    // Получаем токен из localStorage
    const token = localStorage.getItem('token');

    // Если токен есть, отправляем его в заголовках
    if (token) {
      fetch('http://localhost:8080/api/tasks', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`, // Добавляем токен в заголовок
          'Content-Type': 'application/json',
        },
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Failed to fetch tasks');
          }
          return response.json();
        })
        .then(data => setTasks(data))
        .catch(err => setError(err.message));
    } else {
      setError('No token found. Please log in again.');
    }
  }, []);

  return (
    <div>
      <h1>Task List</h1>
      {error && <div className="error">{error}</div>}
      <ul>
        {tasks.map(task => (
          <li key={task.id}>{task.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default TaskList;
