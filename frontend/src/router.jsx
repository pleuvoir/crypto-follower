import { createBrowserRouter } from "react-router-dom";
export default createBrowserRouter(
  [
    {
      path: "/",
      element: <div>/path.</div>,
      errorElement: <div>error 404.</div>,
      children: [
        {
          path: "preferences",
          element: <div>preferences.</div>,
        },
      ],
    },
  ],
  {
    basename: "/app",
  }
);
