import { lazy } from "react";

import { RouteObject } from "react-router-dom";

import Loadable from "../components/third-patry/Loadable";

import FullLayout from "../layout/FullLayout";

const MainPages = Loadable(lazy(() => import("../pages/authentication/Login")));

const Dashboard = Loadable(lazy(() => import("../pages/dashboard")));

const Customer = Loadable(lazy(() => import("../pages/customer")));

const CreateCustomer = Loadable(lazy(() => import("../pages/customer/create")));

const EditCustomer = Loadable(lazy(() => import("../pages/customer/edit")));

const Candidates = Loadable(lazy(() => import("../pages/candidate")));

const CandidateCreate = Loadable(lazy(() => import("../pages/candidate/create")));

const Elections = Loadable(lazy(() => import("../pages/election")));

const VotePage = Loadable(lazy(() => import("../pages/vote")));


const AdminRoutes = (isLoggedIn: boolean): RouteObject => {

  return {

    path: "/",

    element: isLoggedIn ? <FullLayout /> : <MainPages />,

    children: [

      {

        path: "/",

        element: <Dashboard />,

      },

      {

        path: "/candidate",

        children: [

          {

            path: "/candidate",

            element: <Candidates />,

          },

          {

            path: "/candidate/create",

            element: <CandidateCreate />,

          },


        ],

      },

      {

        path: "/election",

        children: [

          {

            path: "/election",

            element: <Elections />,

          },

          {

            path: "/election/:id",

            element: <VotePage />,

          },


        ],

      },


      {

        path: "/customer",

        children: [

          {

            path: "/customer",

            element: <Customer />,

          },

          {

            path: "/customer/create",

            element: <CreateCustomer />,

          },

          {

            path: "/customer/edit/:id",

            element: <EditCustomer />,

          },

        ],

      },

    ],

  };

};


export default AdminRoutes;