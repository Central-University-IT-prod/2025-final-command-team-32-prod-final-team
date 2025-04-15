import * as React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

const AdminPanel = () => {
  return (
    <div className="flex items-center justify-center p-4">
      <div className="w-full max-w-3xl rounded-lg p-8 shadow-xl">
        <h1 className="mb-4 text-center text-4xl font-bold">Admin Panel</h1>
        <p className="mb-8 text-center text-lg">
          Welcome to the Admin Panel. Here you can manage film records.
          <br />
          <span className="text-sm">
            Use the options below to edit, create, or delete films. This panel
            is for administrators only.
          </span>
        </p>
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          <Link to="/adminEdit">
            <Button variant="default" className="w-full py-4 text-lg">
              Edit Film
            </Button>
          </Link>
          <Link to="/adminCreate">
            <Button variant="default" className="w-full py-4 text-lg">
              Create Film
            </Button>
          </Link>
          <Link to="/adminDel">
            <Button variant="destructive" className="w-full py-4 text-lg">
              Delete Film
            </Button>
          </Link>
        </div>
        <div className="mt-8 border-t pt-4">
          <h2 className="mb-2 text-2xl font-semibold">About the Admin Panel</h2>
          <p className="">
            The Admin Panel provides a centralized interface for managing film
            records:
            <br />
            <strong>Edit Film:</strong> Update details of existing films.
            <br />
            <strong>Create Film:</strong> Add new films to the database.
            <br />
            <strong>Delete Film:</strong> Remove films that are no longer needed
            (use with caution).
            <br />
            If you're not familiar with these operations, please consult the
            documentation or contact support.
          </p>
        </div>
      </div>
    </div>
  );
};

export default AdminPanel;
