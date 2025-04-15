import React from "react";
import { useLogin } from "@/Auth";
import { useForm } from "oxyform";
import { Input } from "@/components/ui/lib/input.js";
import { Button } from "@/components/ui/lib/button.js";
import { ErrorWrapper } from "@/ErrorWrapper";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { signUp } from "@/client.js";
import GuestMode from "@/Pages/GuestMode.jsx";
import { CustomAlert } from "@/components/ui/common/custom-alert.jsx";

const Register = () => {
  const form = useForm({
    initialValues: {
      username: "",
      password: "",
      confirmPassword: "",
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
      confirmPassword: [
        "required",
        (x, values) => x === values.password || "Passwords do not match",
      ],
    },
    errors: {
      "regex.username": "test",
      // Дополнительные сообщения об ошибках можно задать здесь
    },
  });

  const login = useLogin();
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const handleFormSubmit = async (values) => {
    try {
      setIsSubmitting(true);
      const response = await signUp(values.username, values.password);
      const tokenUntil = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000); // 7 дней
      login(response.token, tokenUntil, { username: values.username });
      console.log(response.token);
      navigate("/");
      toast.success("Registration successful!");
    } finally {
      setIsSubmitting(false);
    }
  };
  console.log(isSubmitting);

  return (
    <form
      style={{ minWidth: "45vw" }}
      className="mx-auto max-w-xl space-y-4 p-4"
    >
      <div className="space-y-2">
        <label className="block">
          <p className="text-xl"> Username </p>
          <ErrorWrapper {...form.register("username")}>
            <Input
              placeholder="Username"
              className="mt-1 w-full"
              autoComplete="username"
            />
          </ErrorWrapper>
        </label>

        <label className="block">
          <p className="text-xl">Password</p>
          <ErrorWrapper {...form.register("password")}>
            <Input
              type="password"
              placeholder="••••••••"
              className="mt-1 w-full"
              autoComplete="new-password"
            />
          </ErrorWrapper>
        </label>

        <label className="block">
          <p className="text-xl">Confirm Password</p>
          <ErrorWrapper {...form.register("confirmPassword")}>
            <Input
              type="password"
              placeholder="••••••••"
              className="mt-1 w-full"
              autoComplete="new-password"
            />
          </ErrorWrapper>
        </label>
      </div>

      <Button
        type="button"
        onClick={() => form.submit(handleFormSubmit)}
        className="w-full"
        disabled={isSubmitting}
      >
        {isSubmitting ? "Creating account..." : "Sign up"}
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

export default Register;
