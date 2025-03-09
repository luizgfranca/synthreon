import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
    CardFooter,
} from '@/vendor/shadcn/components/ui/card'
import { Label } from '@/vendor/shadcn/components/ui/label'
import { Input } from '@/vendor/shadcn/components/ui/input'
import { Textarea } from '@/vendor/shadcn/components/ui/textarea'
import { Button } from '@/vendor/shadcn/components/ui/button'
import { useForm } from 'react-hook-form'
import { Form, FormField, FormItem, FormMessage } from '@/vendor/shadcn/components/ui/form'
import { useProvider } from '@/context/root'
import { useNavigate, useParams } from 'react-router-dom'
import { ProjectHeader } from '@/component/projectHeader'
import { NewToolDto } from '@/dto/tool.dto'
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod'

const createToolSchema = z.object({
    acronym: z.string()
        .nonempty()
        .regex(/^[a-zA-Z0-9#\$*._-]+$/, {message: 'invalid character(s) in identifier string'})
        .min(3, {message: 'identifier should have at least 3 digits'}),
    name: z.string()
        .nonempty()
        .min(3),
    description: z.string()
})
// TODO: very similar to create project logic
//       there should be a way to generalize this
const CreateToolPage = () => {
    const provider = useProvider()
    const form = useForm<z.infer<typeof createToolSchema>>({
        resolver: zodResolver(createToolSchema)
    })
    const navigate = useNavigate()

    const params = useParams();
    const projectAcronym = params.projectAcronym;
    
    const project = provider.getProjects()
        .find(project => project.acronym === projectAcronym);

    if(!project) {
        // TODO: evaluate showing something here to show to the user
        //       what happened, but in normal usage it shouldn't
        navigate(`${import.meta.env.PL_PATH_PREFIX}/`)
    }

    console.debug('project', project)

    const navigateToProjectHome = () => {
        navigate(`${import.meta.env.PL_PATH_PREFIX}/project/${projectAcronym}`)
    }

    const handleSubmit = (e: unknown) => {
        // TODO: add type validation here for the type conversion
        const data = e as NewToolDto

        if(!projectAcronym) {
            throw new Error('unexpected: trying to create tool but no projectAcronym to be found')
        }

        provider.createTool(projectAcronym, data)
            .then(navigateToProjectHome)
            // FIXME: add proper error handling here
            .catch((e) => console.error('unable to create tool', e))
    }

    // TODO: create a view to encapsulate content related to a project, with the header
    //       and the logic to get the project data
    return (
        <>
            <ProjectHeader projectName={project?.name ?? ''} onLogoutClick={() => {}}/>
            <div className="min-h-screen p-16 bg-zinc-900">
                <div className="max-w-2xl mx-auto">
                    <Card className="bg-zinc-950">
                        <CardHeader>
                            <CardTitle className="text-2xl font-bold">
                                Create New Tool
                            </CardTitle>
                        </CardHeader>
                        <Form {...form}>
                            <form onSubmit={form.handleSubmit(handleSubmit)}>
                                <CardContent className="space-y-6">
                                <div className="space-y-2">
                                    <FormField
                                            control={form.control}
                                            name="acronym"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <Label
                                                        htmlFor="acronym"
                                                        className="text-white"
                                                    >
                                                        Tool Identifiier
                                                    </Label>
                                                    <Input
                                                        placeholder="Enter project identifier"
                                                        className="w-full border-zinc-200 focus:ring-zinc-400"
                                                        {...field}
                                                    />
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    <div className="space-y-2">
                                        <FormField
                                            control={form.control}
                                            name="name"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <Label
                                                        htmlFor="name"
                                                        className="text-white"
                                                    >
                                                        Tool Name
                                                    </Label>
                                                    <Input
                                                        placeholder="Enter project name"
                                                        className="w-full border-zinc-200 focus:ring-zinc-400"
                                                        {...field}
                                                    />
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    <div className="space-y-2">
                                        <FormField
                                            control={form.control}
                                            name="description"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <Label
                                                        htmlFor="description"
                                                        className="text-white"
                                                    >
                                                        Description
                                                    </Label>
                                                    <Textarea
                                                        placeholder="Add a project description"
                                                        className="w-full min-h-32 border-zinc-200 focus:ring-zinc-400"
                                                        {...field}
                                                    />
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </CardContent>
                                <CardFooter className="flex justify-end space-x-4">
                                    <Button
                                        type="button"
                                        variant="outline"
                                        className=""
                                        onClick={navigateToProjectHome}
                                    >
                                        Cancel
                                    </Button>
                                    <Button type="submit" className="">
                                        Create Tool
                                    </Button>
                                </CardFooter>
                            </form>
                        </Form>
                    </Card>
                </div>
            </div>
        </>
    )
}

export default CreateToolPage
