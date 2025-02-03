import { Button } from "@/vendor/shadcn/components/ui/button"
import { Card, CardContent } from "@/vendor/shadcn/components/ui/card"
import { Form, FormControl, FormField, FormItem } from "@/vendor/shadcn/components/ui/form"
import { Input } from "@/vendor/shadcn/components/ui/input"
import { Label } from "@/vendor/shadcn/components/ui/label"
import { useForm } from "react-hook-form"


type PromptProps = {
    title: string,
    onSubmit: (value: string) => void
}

export function Prompt(props: PromptProps) {
    const form = useForm();
    
    return (
        <Card className="space-y-8">
            <CardContent className="py-5 px-10 bg-zinc-850">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit((e) => props.onSubmit(e.prompt))} className="space-y-8">
                        <FormField
                            control={form.control}
                            name="prompt"
                            render={({field}) => (
                                <FormItem>
                                    <Label className="text-lg">{props.title}</Label>
                                    <FormControl>
                                        <Input {...field} />
                                    </FormControl>
                                </FormItem>
                            )}
                        />
                        <Button type="submit">Submit</Button>
                    </form>
                </Form>
                
             </CardContent>
         </Card>
    )
}