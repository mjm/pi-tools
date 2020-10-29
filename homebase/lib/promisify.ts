export type PromisifiedClient<T> = {
    [P in keyof T]:
    T[P] extends (req: infer Req, cb: (err: any, res: infer Res | null) => void) => any
        ? (p1: Req) => Promise<Res> :
        never;
}

export function promisifyClient<T>(client: T): PromisifiedClient<T> {
    const promiseClient: PromisifiedClient<T> = {} as any;
    for (const key of Object.keys(Object.getPrototypeOf(client))) {
        promiseClient[key] = promisify(client, key);
    }
    return promiseClient;
}

export function promisify(obj, methodName): Function {
    return (...args) => {
        return new Promise((resolve, reject) => {
            args.push((err, res) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(res);
                }
            });

            obj[methodName](...args);
        });
    };
}
