import os
from pprint import pprint

PATH = os.getenv("GOPATH") + '/src/github.com/trevorprater/youfie-api-2/routers/'
for f in os.listdir(PATH):
    routes = []
    methods = []
    with open(PATH + f) as ff:
        data = []
        route = methods = ''
        for line in ff.readlines():
            if 'router.Handle(' in line or 'router.HandleFunc(' in line:
                route = line[line.find('(') + 1:line.find(',')].replace(
                    '"', '').strip()
            if 'Methods(' in line:
                methods = line.split('Methods(')[-1].replace(')', '').replace(
                    '"', '').strip()
            if not route == '' and not methods == '':
                data.append('endpoint: {},  methods: {}'.format(route, methods))
                route = methods = ''
        if data:
            print f
            print '=' * 50
            for row in data:
                print row
    print '\n'
