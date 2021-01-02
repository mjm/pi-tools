export function destinationURL(shortURL: string): string {
    if (process.env.NODE_ENV === "development") {
        return `http://localhost:4240/${shortURL}`;
    } else {
        return `http://go/${shortURL}`;
    }
}
