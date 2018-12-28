import pdb

def trace(errors=None, on_input=None, on_output=None):
    """
    trace() returns a decorator that will drop the file into a pdb session if
    the appropriate condition is met. Samle usage

    >>> @trace(errors=Exception)
    >>> def foo():
    >>>     ...
    >>> foo()  # opens pdb if foo() throws an error

    >>> @trace(on_input=lambda arg: arg == 1)
    >>> def foo(arg):
    >>>     ...
    >>> foo(1)  # opens pdb if foo() gets an input of 1

    >>> @trace(on_output=lambda ret: ret == 1)
    >>> def foo():
    >>>     ...
    >>> foo()  # opens pdb if foo() returns 1

    on_input: (*any) => bool
    errors: (Exception,)
    on_output: (*any) => bool
    """

    def decorator(f):
        def g(*args, **kwargs):
            if on_input is not None and on_input(*args, **kwargs):
                pdb.set_trace()

            if errors is not None:
                try:
                    retval = f(*args, **kwargs)
                except errors as e:
                    pdb.set_trace()
            else:
                retval = f(*args, **kwargs)

            if on_output is not None and on_output(retval):
                pdb.set_trace()

            return retval
        return g
    return decorator

