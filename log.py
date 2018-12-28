records =  []
indent = 0

def pretty_print_records(records):
    for record in records:
        indent='  '*record['indent']
        fname=record['fname']
        retval=record['retval']
        args=', '.join(map(str, record['args']))
        kwargs=', '.join(
            '{}={}'.format(argname, argval)
            for argname, argval in record['kwargs'].iteritems()
        )
        print('{indent}{fname}({args}{kwargs}) -> {retval}'.format(
            indent=indent,
            fname=fname,
            args=args,
            kwargs=kwargs,
            retval=retval
        ))

def log(f):
    def g(*args, **kwargs):
        global indent, records
        this_record = {
            "indent": indent,
            "args": args,
            "kwargs": kwargs,
            "fname": f.__name__,
        }
        records.append(this_record)
        indent += 1
        retval = f(*args, **kwargs)
        indent -= 1
        this_record["retval"] = retval
        if indent == 0:
            pretty_print_records(records)
        return retval
    return g

@log
def fib(n):
    if n <= 2: return 1
    return fib(n - 1) + fib(n - 2)

fib(4)
print('===============================')


@log
def h(n):
    return n**2

@log
def g(n):
    return h(3*n)

@log
def f(n):
    return g(n + 1)


f(5)
