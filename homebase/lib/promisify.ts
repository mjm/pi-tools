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
