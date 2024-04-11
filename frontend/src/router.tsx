import { createBrowserRouter } from 'react-router-dom';
import MainLayout from './Layout';
import React from 'react';
import Tasks from './View/Page/Tasks';
import Preferences from './View/Page/Preferences';
import Document from './View/Page/Document';

export default createBrowserRouter(
  [
    {
      path: '/',
      element: <MainLayout />,
      errorElement: <div>error 404.</div>,
      children: [
        {
          path: '/',
          element: <Document />,
        },
        {
          path: 'tasks',
          element: <Tasks />,
        },
        {
          path: 'preferences',
          element: <Preferences />,
        },
      ],
    },
  ],
  {
    basename: '/app',
  },
);
