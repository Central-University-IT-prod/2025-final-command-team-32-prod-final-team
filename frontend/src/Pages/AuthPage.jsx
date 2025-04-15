import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/lib/tabs";
import Login from "@/Pages/Login.jsx";
import Register from "@/Pages/Register.jsx";
import React from "react";

export function AuthPage() {
  return (
    <div style={{ minWidth: "45vw" }} className="mx-auto w-full max-w-xl">
      <Tabs defaultValue="login">
        <div>
          <TabsList
            style={{ minWidth: "45vw" }}
            className="mx-auto grid w-full max-w-xl grid-cols-2"
          >
            <TabsTrigger value="login">Login</TabsTrigger>
            <TabsTrigger value="register">Register</TabsTrigger>
          </TabsList>
          <TabsContent value="login">
            <Login />
          </TabsContent>
          <TabsContent value="register">
            <Register />
          </TabsContent>
        </div>
      </Tabs>
    </div>
  );
}
