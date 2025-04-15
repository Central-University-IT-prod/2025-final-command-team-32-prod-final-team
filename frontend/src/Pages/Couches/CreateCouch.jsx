import React from "react";
import { Input } from "@/components/ui/lib/input.js";
import { ErrorWrapper } from "@/ErrorWrapper.jsx";
import { useForm } from "oxyform";
import { Button } from "@/components/ui/lib/button.js";
import { createCouch, searchUsers } from "@/client.js";
import { useAuth } from "@/Auth.js";
import { toast } from "sonner";
import MultipleSelector from "@/components/ui/common/multiple-selector.js";
import { useNavigate } from "react-router-dom";

const CreateCouch = () => {
  const navigate = useNavigate();
  const [, , , token] = useAuth();
  const [userQuery, setUserQuery] = React.useState("");
  const [foundUsers, setFoundUsers] = React.useState([]);
  const [cachedUsers, setCachedUsers] = React.useState([]);

  const form = useForm({
    initialValues: {
      title: "",
      users: [],
    },
    validation: {
      title: ["required", /^[a-zA-Zа-яА-Я0-9.,!?:;()\s"'«»—–-]+$/],
      users: [
        (value) =>
          (Array.isArray(value) && value.length > 0) ||
          "Select at least one user",
      ],
    },
  });

  React.useEffect(() => {
    if (userQuery.length > 2) {
      searchUsers(userQuery)
        .then((users) => {
          setFoundUsers(users);
          setCachedUsers((prev) => [
            ...prev,
            ...users.filter((u) => !prev.some((pu) => pu.id === u.id)),
          ]);
        })
        .catch(() => {
          toast.error("Error fetching users");
        });
    } else {
      setFoundUsers([]);
    }
  }, [userQuery]);

  // Исправляем преобразование с проверкой на undefined
  const selectedOptions = (form.values.users || []).map((userId) => {
    const user = cachedUsers.find((u) => u.id === userId);
    return {
      value: userId,
      label: user?.username || "Unknown User",
    };
  });

  const handleSelectorChange = (newOptions) => {
    const userIds = newOptions.map((opt) => opt.value);
    form.setValue("users", userIds);
  };

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const handleSubmit = async (values) => {
    try {
      setIsSubmitting(true);

      if (!values.users || values.users.length === 0) {
        toast.error("Please select at least one user");
        return;
      }

      const selectedUsernames = values.users
        .map((userId) => {
          const user = cachedUsers.find((u) => u.id === userId);
          return user?.username;
        })
        .filter(Boolean); // Убираем undefined, если user не найден

      const response = await createCouch(
        values.title,
        selectedUsernames,
        token,
      );
      toast.success("Couch successfully created!");
      navigate(`/couches/${response.id}`);
    } catch (error) {
      console.error(error);
      toast.error("Failed to create couch. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      style={{ minWidth: "45vw" }}
      className="mx-auto max-w-xl space-y-4 p-4"
      onSubmit={(e) => {
        e.preventDefault();
        form.submit(handleSubmit);
      }}
    >
      <div className="space-y-2">
        <h1 style={{ fontSize: "32px", marginBottom: "16px" }}>Create Couch</h1>
        <label className="block">
          <p className="text-xl">Couch name</p>
          <ErrorWrapper {...form.register("title")}>
            <Input
              placeholder="Enter couch name"
              className="mt-2 w-full"
              disabled={isSubmitting}
            />
          </ErrorWrapper>
        </label>
        <div>
          <p className="text-xl">Select users</p>
          <MultipleSelector
            disabled={isSubmitting}
            placeholder="Start typing a username to add"
            options={foundUsers.map((user) => ({
              value: user.id,
              label: user.username,
            }))}
            value={selectedOptions}
            onValueChange={handleSelectorChange}
            inputValue={userQuery}
            onInputChange={setUserQuery}
          />
          {form.errors.users && (
            <div className="text-destructive mt-1 text-sm">
              {form.errors.users}
            </div>
          )}
        </div>
      </div>
      <Button type="submit" className="w-full text-lg" disabled={isSubmitting}>
        {isSubmitting ? "Creating..." : "Create couch"}
      </Button>
    </form>
  );
};

export default CreateCouch;
