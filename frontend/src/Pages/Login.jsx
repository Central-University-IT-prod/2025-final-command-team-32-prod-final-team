import React from "react";
import { useLogin } from "@/Auth";
import { useForm } from "oxyform";
import { Input } from "@/components/ui/lib/input.js";
import { Button } from "@/components/ui/lib/button.js";
import { ErrorWrapper } from "@/ErrorWrapper";
import { useNavigate } from "react-router";
import { signIn } from "@/client.js";
import GuestMode from "@/Pages/GuestMode.jsx";
import { CustomAlert } from "@/components/ui/common/custom-alert.jsx";

const Login = () => {
  const form = useForm({
    initialValues: {
      username: "",
      password: "",
    },
    validation: {
      username: [
        "required",
        (x) => x.length >= 3 || "Username must be at least 3 characters",
        /^[a-zA-Zа-яА-Я0-9.,!?:;()\s"'«»—–-]+$/,
      ],
      password: [
        "required",
        (x) => x.length >= 8 || "Password must be at least 8 characters long",
      ],
    },
    errors: {
      required: "This field is required",
      "regex.username":
        "Username can contain only English letters, numbers, dots, and parentheses",
    },
  });

  const login = useLogin();
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const handleSubmit = async (values) => {
    try {
      setIsSubmitting(true);
      const response = await signIn(values.username, values.password);
      const tokenUntil = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000);
      login(response.token, tokenUntil, { username: values.username });

      const savedUrl = localStorage.getItem("refer_from");
      if (savedUrl) {
        navigate(savedUrl, { replace: true });
        localStorage.removeItem("refer_from");
      } else {
        navigate("/");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      style={{ minWidth: "45vw" }}
      className="mx-auto max-w-xl space-y-4 p-4"
    >
      <div className="space-y-2">
        <label className="block">
          <p className="text-xl">Username</p>
          <ErrorWrapper {...form.register("username")}>
            <Input
              placeholder="Username"
              className="mt-2 w-full"
              autoComplete="username"
              disabled={isSubmitting}
            />
          </ErrorWrapper>
        </label>
        <label className="block">
          <p className="text-xl">Password</p>
          <ErrorWrapper {...form.register("password")}>
            <Input
              type="password"
              placeholder="••••••••"
              className="mt-2 w-full"
              autoComplete="new-password"
              disabled={isSubmitting}
            />
          </ErrorWrapper>
        </label>
      </div>
      <Button
        type="button"
        onClick={() => form.submit(handleSubmit)}
        className="w-full text-lg"
        disabled={isSubmitting}
      >
        {isSubmitting ? "Signing in..." : "Sign in"}
      </Button>
      <h1
        style={{
          textAlign: "center",
          fontSize: "32px",
          fontWeight: "bold",
          marginTop: "2rem",
        }}
      >
        Or
      </h1>
      <CustomAlert title="Dont want to register?" message="Try guest mode!" />
      <GuestMode
        isSubmitting={isSubmitting}
        setIsSubmitting={setIsSubmitting}
      />
    </form>
  );
};

export default Login;
