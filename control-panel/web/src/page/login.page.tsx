import { Button } from "@/vendor/shadcn/components/ui/button";
import { Card, CardContent } from "@/vendor/shadcn/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/vendor/shadcn/components/ui/form";
import { Input } from "@/vendor/shadcn/components/ui/input";
import { Label } from "@/vendor/shadcn/components/ui/label";
import { Separator } from "@radix-ui/react-separator";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

export function LoginPage() {
  const form = useForm();
  const navigate = useNavigate();

  return (
    <div className="flex h-screen items-center justify-center">
      <Card className="space-y-8">
        <CardContent className="py-5 px-10 bg-zinc-850">
          <h1 className="text-2xl font-bold pb-8">Enter your control panel</h1>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit((e) => {
                console.log("onsubmit", e);
                navigate("/");
              })}
              className="space-y-8"
            >
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input type="email" placeholder="E-Mail" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Password"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <div className="flex w-full justify-center">
                <Button type="submit">Submit</Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
