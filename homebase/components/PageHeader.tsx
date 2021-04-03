export default function PageHeader({children, buttons}: {
    children: React.ReactNode;
    buttons?: React.ReactNode;
}) {
    return (
        <header className="bg-white shadow">
            <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                <div className="md:flex md:items-center md:justify-between">
                    <div className="flex-1 min-w-0">
                        <h1 className="text-2xl font-bold leading-tight text-gray-900 sm:text-3xl sm:truncate">
                            {children}
                        </h1>
                    </div>
                    {buttons ? (
                        <div className="mt-4 flex md:mt-0 md:ml-4">
                            {buttons}
                        </div>
                    ) : null}
                </div>
            </div>
        </header>
    );
}
