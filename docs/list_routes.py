import os
from pprint import pprint

PATH = os.getenv("GOPATH") + '/src/github.com/trevorprater/youfie-api-2/routers/'
endpoint_dict = {}
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
                try:
                    endpoint_dict[route] = endpoint_dict[route] + "," + methods
                except KeyError as e:
                    endpoint_dict[route] = methods
                data.append("'{}': {}".format(route, methods))
                route = methods = ''
        #if data:
            #for row in data:
            #    print row
for endpoint  in endpoint_dict.keys():
    print str(endpoint)+ ": " + ", ".join(endpoint_dict[endpoint].split(','))
