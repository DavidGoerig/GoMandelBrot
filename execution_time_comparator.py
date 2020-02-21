
#########################################################
# @author:  David Goerig                                #
# @id:      djg53                                       #
# @module:  Concurrency and Parallelism - CO890         #
# @asses:   assess 4- Go Worker and geometric           #
#           distribution comparaison                    #
#########################################################
###
# desc: compare the fastest way to read a large file (nested loop, classic object programming, coroutines), and show the result in a file
###
import os
import argparse as parse
import time
import math

##
# param:    /
# desc:     management of the argument of the script(-f, -o, -n, -d). argparse is used
##


def manage_args():
    outputfilename = "compare_result.txt"
    iteration = 100
    execoutput = ".outputs_to_delete.txt"
    parser = parse.ArgumentParser()
    # creation of the arguments, with the arg type, flag name and text in the -h
    parser.add_argument("-o", "-outputfilename", type=int, help="outputfilename (default: compare_result.txt)")
    parser.add_argument("-n", "-nbrit", type=int, help="nbr of iteration (default: 100)")
    parser.add_argument("-d", "-execoutput", type=str, help="file for the output, None for stantard output (default: .outputs_to_delete.txt)")
    args = parser.parse_args()
    # assignement of the argument to the good param
    if args.o:
        outputfilename = args.o
    if args.d:
        execoutput = args.d
    if args.n:
        if args.n <= 0:
            print("Iteration nbr need to be a positive integer")
            return 1
        iteration = args.n
    return iteration, execoutput, outputfilename

##
# param:    dir_path: path of the current working directory, params: script parameter (here the filename), execoutput: redirection file
# desc:     exec the python scripts with the file as argument and redirect in output if wanted
##


def build_script(dir_path, go_scrit, execoutput):
    if execoutput == "None":
        execoutput = ""
    else:
        execoutput = " > " + execoutput
    os.system("go build  " + dir_path + "/" + go_scrit + " " + execoutput)
##
# param:    dir_path: path of the current working directory, params: script parameter (here the filename), execoutput: redirection file
# desc:     exec the python scripts with the file as argument and redirect in output if wanted
##


def exec_script(dir_path, go_scrit, execoutput, arg):
    if execoutput == "None":
        execoutput = ""
    else:
        execoutput = " > " + execoutput
    os.system("." + "/" + go_scrit + " " + str(arg) + " " + execoutput)
##
# param:    inputfilename: name of the file to read, iteration: nbr of iteration, execoutput: redirection file
# desc:     loop on each scripts (in python_file_name), execute them, and return execution time in a dictionnary
##


def exec_them(iteration, execoutput):
    go_file_name = {
        "geo_distrib": "Graphical_geo_distrib.go",
        "workers": "Graphical_worker.go"
    }
    go_experiment = {
        "geo_distrib_divide_in_5": "Graphical_geo_distrib",
        "geo_distrib_divide_in_10": "Graphical_geo_distrib",
        "geo_distrib_divide_in_20": "Graphical_geo_distrib",
        "geo_distrib_divide_in_50": "Graphical_geo_distrib",
        "geo_distrib_divide_in_100": "Graphical_geo_distrib ",
        "worker_with_5_workers": "Graphical_worker",
        "worker_with_10_workers": "Graphical_worker",
        "worker_with_20_workers": "Graphical_worker",
        "worker_with_50_workers": "Graphical_worker",
        "worker_with_100_workers": "Graphical_worker",
    }
    go_nbr = {
        "geo_distrib_divide_in_5": 5,
        "geo_distrib_divide_in_10": 10,
        "geo_distrib_divide_in_20": 20,
        "geo_distrib_divide_in_50": 50,
        "geo_distrib_divide_in_100": 100,
        "worker_with_5_workers": 5,
        "worker_with_10_workers": 10,
        "worker_with_20_workers": 20,
        "worker_with_50_workers": 50,
        "worker_with_100_workers": 100,
    }
    result = {
        "geo_distrib_divide_in_5": [],
        "geo_distrib_divide_in_10": [],
        "geo_distrib_divide_in_20": [],
        "geo_distrib_divide_in_50": [],
        "geo_distrib_divide_in_100": [],
        "worker_with_5_workers": [],
        "worker_with_10_workers": [],
        "worker_with_20_workers": [],
        "worker_with_50_workers": [],
        "worker_with_100_workers": [],
    }
    dir_path = os.path.dirname(os.path.realpath(__file__))
    for x in go_file_name:
        build_script(dir_path, go_file_name[x], execoutput)
    for x in go_experiment:
        for it in range(0, iteration):
            # start time counter
            start_time = time.time()
            # exec scripts
            exec_script(dir_path, go_experiment[x], execoutput, go_nbr[x])
            # calc execution time
            result[x].append(time.time() - start_time)
    return result
##
# param:    results: time of work from each iteration of each scripts in a dictionary
# desc:     calc for each script: the  average,  standard  deviation,  and  the standard error of the mean
##


def compute_result(results):
    # dictionaries containing results for each scripts
    average = {}
    sd = {}
    sem = {}
    # loop in each script, and calc average, sd, sem and assign them to the good dictionary
    for x in results:
        total = 0
        for i in results[x]:
            total += i
        # average calc
        loc_average = total / len(results[x])
        sum_mean_diff = 0
        for i in results[x]:
            sum_mean_diff += pow((i - loc_average), 2)
        # standard deviation calc
        loc_sd = math.sqrt(sum_mean_diff / len(results[x]))
        # stand error to the mean calc
        loc_sem = loc_sd / math.sqrt(len(results[x]))
        # assignation to the wanted script in the dictionary
        average[x] = loc_average
        sd[x] = loc_sd
        sem[x] = loc_sem
    return average, sd, sem


##
# param:    dictionary_to_sort
# desc:     find a min in a dict and return the key
##
def find_min(dictionary_to_sort):
    res = -1
    val = -1
    for x in dictionary_to_sort:
        if res == -1:
            res = x
            val = dictionary_to_sort[x]
        if dictionary_to_sort[x] < val:
            res = x
            val = dictionary_to_sort[x]
    return res
##
# param:    dictionary_to_sort_source
# desc:     create an array with the key in the wanted order
##

def sort_dict_min_to_max(dictionary_to_sort_source):
    dictionary_to_sort = dict(dictionary_to_sort_source)
    order = []
    while len(dictionary_to_sort) != 0:
        key = find_min(dictionary_to_sort)
        order.append(key)
        del dictionary_to_sort[key]
    return order

##
# param:    outputfilename: file name for the results, average: average for each scripts,
#           sd: standard deviation for each scripts, sem:  for each
# desc:     create the file with information for each scripts, then order them and show the % diff
##


def create_file(outputfilename, average, sd, sem, iteration):
    content = ""
    for x in average:
        content += "---------  " + x + " ---------" + "\n"
        content += "Iteration nbr:\t\t\t\t\t" + str(iteration) + "\n"
        content += "Average:\t\t\t\t\t\t" + str(average[x]) + "\n"
        content += "Standard deviation:\t\t\t\t" + str(sd[x]) + "\n"
        content += "Standard error of the mean:\t\t" + str(sem[x]) + "\n" + "\n"
    content += "Fastest method in order according to the average:\n"

    # sort averages, calc % diff and print
    # TODO: ICI JUSTE SORT CA, ENSUITE L'APPLIQUER A SD ET SEM
    sorted_average = sort_dict_min_to_max(average)
    last = 0
    for x in sorted_average:
        diff = ""
        if (last != 0):
            diff = " (" + str(((average[x] - last) / last) * 100) + " % more than the previous one)."
        last = average[x]
        content += x + "\t->\t" + str(average[x]) + diff + "\n"
    content += "\nLess deviation from the average (standard deviation):\n"
    # sort standard deviation, calc % diff and print
    sorted_sd = sort_dict_min_to_max(sd)
    last = 0
    for x in sorted_sd:
        diff = ""
        if (last != 0):
            diff = " (" + str(((sd[x] - last) / last) * 100) + " % more than the previous one)."
        last = sd[x]
        content += x + "\t->\t" + str(sd[x]) + diff +  "\n"
    content += "\nMost homogeneous method in order (standard error to the mean):\n"
    # sort standard error to the, calc % diff and print
    sorted_sem = sort_dict_min_to_max(sem)
    last = 0
    for x in sorted_sem:
        diff = ""
        if last != 0:
            diff = " (" + str(((sem[x] - last) / last) * 100) + " % more than the previous one)."
        last = sem[x]
        content += x + "\t->\t" + str(sem[x]) + diff +  "\n"
    # open the file, overwrite it, write in it, and then close the file
    f = open(outputfilename, "w")
    f.write(content)
    f.close()
##
# param:    /
# desc:     entry point, call all the function in order: check the parameters, test the file,
#           execute script, calc of the wanted value (sd, mean, sem), file creation
##


def main():
    iteration, execoutput, outputfilename = manage_args()
    results = exec_them(iteration, execoutput)
    average, sd, sem = compute_result(results)
    create_file(outputfilename, average, sd, sem, iteration)
    return 0
##
# param:    /
# desc:     entry point
##


if __name__ == '__main__':
    exit(main())
