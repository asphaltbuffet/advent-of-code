import datetime
from cookiecutter.utils import simple_filter


@simple_filter
def camel_case(v):
    res = v.split(" ")
    for i in range(len(res)):
        f = lambda x: x.upper()
        if i == 0:
            f = lambda x: x.lower()

        res[i] = f(res[i][0]) + res[i][1:]

    return "".join(res)


@simple_filter
def current_year(_):
    return datetime.datetime.now().year


@simple_filter
def current_day(_):
    return datetime.datetime.now().day