type LoadingStatus = 'none' | 'pending' | 'success' | 'fail';

export class Datasource<T> {
    #loadingStatus: LoadingStatus = 'none';
    #loadingPromise?: Promise<T>;
    
    #data?: T;
    #error: unknown;

    async fetch(): Promise<T> {
        throw new Error('fetch() should be overriden with the function to load your dataset')
    }

    get(): T {
        switch(this.#loadingStatus) {
            case 'none':
                this.#loadingPromise = this.fetch();
                this.#loadingPromise.
                    then((data) => {
                        this.#data = data;
                        this.#loadingStatus = 'success'
                    }).catch(e => {
                        this.#error = e;
                        this.#loadingStatus = 'fail'
                    });
                this.#loadingStatus = 'pending'
                throw this.#loadingPromise;
            case 'pending':
                throw this.#loadingPromise;
            case 'fail':
                throw this.#error
            case 'success':
                return this.#data as T;
        }
    }
}