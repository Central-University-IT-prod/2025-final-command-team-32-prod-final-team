import React, { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import "./index.css";
import { AuthProvider } from "./Auth.tsx";
import { Cinema } from "@/Pages/Cinema.jsx";
import "./font.css";
import "./button.css";
import "./overflows.css";
import { Navbar } from "@/Navbar.jsx";
import { Toaster } from "sonner";
import MyLikes from "@/Pages/MyLikes/MyLikes.jsx";
import { Explore } from "@/Pages/Explore/Explore.jsx";
import { AutoAlt } from "@/AutoAlt.jsx";
import "bootstrap-icons/font/bootstrap-icons.min.css";
import { Search } from "@/Pages/Search.jsx";
import { Popular } from "@/Pages/Popular.jsx";
import { CreateAccountDialogProvider } from "@/CreateAccountDialog.jsx";
import Create from "@/Pages/Create.jsx";
import { Couches } from "@/Pages/Couches/Couches.jsx";
import CreateCouch from "@/Pages/Couches/CreateCouch.jsx";
import { Couche } from "@/Pages/Couches/Couche.jsx";
import { CoucheFeed } from "@/Pages/Couches/CoucheFeed.jsx";
import CouchLikes from "@/Pages/Couches/CouchLikes.jsx";
import AdminCreate from "@/Pages/Admin/AdminCreate.jsx";
import AdminPanel from "@/Pages/Admin/AdminPanel.jsx";
import { AuthPage } from "@/Pages/AuthPage.jsx";
import { AdminRouteGuard } from "@/AdminRouteGuard.js";
import AdminEditFilm from "@/Pages/Admin/AdminEditFilm.js";
import AdminDelFilm from "@/Pages/Admin/AdminDelFilm.js";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <AuthProvider config={{ storeName: "auth", headerFunction: (x) => ({}) }}>
      <BrowserRouter>
        <CreateAccountDialogProvider>
          <Navbar />
          <div className="h-[32px] md:h-[128px]" />
          <Routes>
            <Route
              path={"/"}
              element={<Navigate to="/popular" replace={true} />}
            />
            <Route
              path={"/login"}
              element={<Navigate to="/auth" replace={true} />}
            />
            <Route
              path={"/register"}
              element={<Navigate to="/auth" replace={true} />}
            />
            <Route path={"/couches/:id"} element={<Couche />} />
            <Route path={"/couches/:id/feed"} element={<CoucheFeed />} />
            <Route path={"/couches/:id/likes"} element={<CouchLikes />} />
            <Route path={"/search"} element={<Search />} />
            <Route path={"/couches"} element={<Couches />} />
            <Route path={"/popular"} element={<Popular />} />
            <Route path={"/explore"} element={<Explore />} />
            <Route path={"/cinema/:id"} element={<Cinema />} />
            <Route path={"/myLikes"} element={<MyLikes />} />
            <Route path={"/create"} element={<Create />} />
            <Route path={"/createCouch"} element={<CreateCouch />} />
            <Route path={"/auth"} element={<AuthPage />} />
            <Route
              path="/adminPanel"
              element={
                <AdminRouteGuard>
                  <AdminPanel />
                </AdminRouteGuard>
              }
            />
            <Route
              path="/adminCreate"
              element={
                <AdminRouteGuard>
                  <AdminCreate />
                </AdminRouteGuard>
              }
            />
            <Route
              path="/adminEdit"
              element={
                <AdminRouteGuard>
                  <AdminEditFilm />
                </AdminRouteGuard>
              }
            ></Route>
            <Route
              path="/adminDel"
              element={
                <AdminRouteGuard>
                  <AdminDelFilm />
                </AdminRouteGuard>
              }
            ></Route>
          </Routes>
          <div className="h-[72px] md:h-[128px]" />
          <Toaster />
        </CreateAccountDialogProvider>
      </BrowserRouter>
    </AuthProvider>
    <AutoAlt />
  </StrictMode>,
);

const root = window.document.documentElement;
root.classList.remove("light", "dark");
root.classList.add("dark");
