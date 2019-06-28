import datetime
import random
import sys

ProductCode = {
    'Xmanager': 0,
    'Xshell': 1,
    'Xlpd': 2,
    'Xfile': 3,
    'Xftp': 4,
    'Xmanager 3D': 5,
    'Xmanager Enterprise': 6,
    'Xshell Plus': 7
}

LicenseType = [
    [ProductCode['Xmanager'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xmanager'], 0x0C, 0, 'Educational', 2],
    [ProductCode['Xmanager'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xmanager'], 0x10, 0, 'Educational', 1],
    [ProductCode['Xmanager'], 0x16, 2, 'Student 2-year Subscription', 2],
    [ProductCode['Xmanager'], 0x18, 4, 'Student 4-year Subscription', 2],
    [ProductCode['Xmanager'], 0x20, 2, 'Student 2-year Subscription', 1],
    [ProductCode['Xmanager'], 0x22, 4, 'Student 4-year Subscription', 1],
    [ProductCode['Xmanager'], 0x3D, 0, 'Standard Subscription', 2],
    [ProductCode['Xmanager'], 0x3E, 0, 'Educational Subscription', 2],
    [ProductCode['Xmanager'], 0x41, 0, 'Standard Subscription', 1],
    [ProductCode['Xmanager'], 0x42, 0, 'Educational Subscription', 1],
    [ProductCode['Xmanager'], 0x47, 0, 'Standard Subscription', 2],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x48, 0, 'Educational Subscription', 2],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x4B, 0, 'Standard Subscription', 1],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x4C, 0, 'Educational Subscription', 1],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x51, 0, 'Standard', 2],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x52, 0, 'Educational', 2],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x55, 0, 'Standard', 1],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x56, 0, 'Educational', 1],  # Concurrent Registered
    [ProductCode['Xmanager'], 0x60, 0, 'Standard', 1],
    [ProductCode['Xmanager'], 0x61, 0, 'Standard', 2],
    [ProductCode['Xmanager'], 0x62, 0, 'Standard', 1],
    [ProductCode['Xmanager'], 0x63, 0, 'Standard', 2],
    [ProductCode['Xmanager'], 0x29, 1, 'CLS Class A', 2],
    [ProductCode['Xmanager'], 0x2A, 1, 'CLS Class B', 2],
    [ProductCode['Xmanager'], 0x2B, 1, 'CLS Class C', 2],
    [ProductCode['Xmanager'], 0x2C, 1, 'DLS', 2],
    [ProductCode['Xmanager'], 0x2D, 1, 'SLS', 2],
    [ProductCode['Xmanager'], 0x33, 1, 'CLS Class A', 1],
    [ProductCode['Xmanager'], 0x34, 1, 'CLS Class B', 1],
    [ProductCode['Xmanager'], 0x35, 1, 'CLS Class C', 1],
    [ProductCode['Xmanager'], 0x36, 1, 'DLS', 1],
    [ProductCode['Xmanager'], 0x37, 1, 'SLS', 1],
    [ProductCode['Xshell Plus'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xshell'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xshell'], 0x0C, 0, 'Educational', 2],
    [ProductCode['Xshell'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xshell'], 0x10, 0, 'Educational', 1],
    [ProductCode['Xshell'], 0x16, 2, 'Student 2-year Subscription', 2],
    [ProductCode['Xshell'], 0x18, 4, 'Student 4-year Subscription', 2],
    [ProductCode['Xshell'], 0x20, 2, 'Student 2-year Subscription', 1],
    [ProductCode['Xshell'], 0x22, 4, 'Student 4-year Subscription', 1],
    [ProductCode['Xshell'], 0x3D, 0, 'Standard Subscription', 2],
    [ProductCode['Xshell'], 0x3E, 0, 'Educational Subscription', 2],
    [ProductCode['Xshell'], 0x41, 0, 'Standard Subscription', 1],
    [ProductCode['Xshell'], 0x42, 0, 'Educational Subscription', 1],
    [ProductCode['Xshell'], 0x47, 0, 'Standard Subscription', 2],
    [ProductCode['Xshell'], 0x48, 0, 'Educational Subscription', 2],
    [ProductCode['Xshell'], 0x4B, 0, 'Standard Subscription', 1],
    [ProductCode['Xshell'], 0x4C, 0, 'Educational Subscription', 1],
    [ProductCode['Xshell'], 0x51, 0, 'Standard', 2],
    [ProductCode['Xshell'], 0x52, 0, 'Educational', 2],
    [ProductCode['Xshell'], 0x55, 0, 'Standard', 1],
    [ProductCode['Xshell'], 0x56, 0, 'Educational', 1],
    [ProductCode['Xshell'], 0x60, 0, 'Standard', 1],  # ������
    [ProductCode['Xshell'], 0x61, 0, 'Standard', 2],  # ������
    [ProductCode['Xshell'], 0x62, 0, 'Standard', 1],
    [ProductCode['Xshell'], 0x63, 0, 'Standard', 2],
    [ProductCode['Xlpd'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xlpd'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xlpd'], 0x3D, 0, 'Standard Subscription', 2],
    [ProductCode['Xlpd'], 0x3E, 0, 'Educational Subscription', 2],
    [ProductCode['Xlpd'], 0x41, 0, 'Standard Subscription', 1],
    [ProductCode['Xlpd'], 0x42, 0, 'Educational Subscription', 1],
    [ProductCode['Xlpd'], 0x47, 0, 'Standard Subscription', 2],
    [ProductCode['Xlpd'], 0x48, 0, 'Educational Subscription', 2],
    [ProductCode['Xlpd'], 0x4B, 0, 'Standard Subscription', 1],
    [ProductCode['Xlpd'], 0x4C, 0, 'Educational Subscription', 1],
    [ProductCode['Xlpd'], 0x51, 0, 'Standard', 2],
    [ProductCode['Xlpd'], 0x55, 0, 'Standard', 1],
    [ProductCode['Xlpd'], 0x60, 0, 'Standard', 1],
    [ProductCode['Xlpd'], 0x61, 0, 'Standard', 2],
    [ProductCode['Xlpd'], 0x62, 0, 'Standard', 1],
    [ProductCode['Xlpd'], 0x63, 0, 'Standard', 2],
    [ProductCode['Xfile'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xftp'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xftp'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xftp'], 0x3D, 0, 'Standard Subscription', 2],
    [ProductCode['Xftp'], 0x3E, 0, 'Educational Subscription', 2],
    [ProductCode['Xftp'], 0x41, 0, 'Standard Subscription', 1],
    [ProductCode['Xftp'], 0x42, 0, 'Educational Subscription', 1],
    [ProductCode['Xftp'], 0x47, 0, 'Standard Subscription', 2],
    [ProductCode['Xftp'], 0x48, 0, 'Educational Subscription', 2],
    [ProductCode['Xftp'], 0x4B, 0, 'Standard Subscription', 1],
    [ProductCode['Xftp'], 0x4C, 0, 'Educational Subscription', 1],
    [ProductCode['Xftp'], 0x51, 0, 'Standard', 2],
    [ProductCode['Xftp'], 0x55, 0, 'Standard', 1],
    [ProductCode['Xftp'], 0x60, 0, 'Standard', 1],
    [ProductCode['Xftp'], 0x61, 0, 'Standard', 2],
    [ProductCode['Xftp'], 0x62, 0, 'Standard', 1],
    [ProductCode['Xftp'], 0x63, 0, 'Standard', 2],
    [ProductCode['Xmanager 3D'], 0x0B, 0, 'Standard', 2],
    [ProductCode['Xmanager 3D'], 0x0C, 0, 'Educational', 2],
    [ProductCode['Xmanager 3D'], 0x0F, 0, 'Standard', 1],
    [ProductCode['Xmanager 3D'], 0x10, 0, 'Educational', 1],
    [ProductCode['Xmanager Enterprise'], 0x0B, 0, '', 2],
    [ProductCode['Xmanager Enterprise'], 0x0C, 0, 'Educational', 2],
    [ProductCode['Xmanager Enterprise'], 0x0F, 0, '', 1],
    [ProductCode['Xmanager Enterprise'], 0x10, 0, 'Educational', 1],
    [ProductCode['Xmanager Enterprise'], 0x3D, 0, 'Standard Subscription', 2],
    [ProductCode['Xmanager Enterprise'], 0x3E, 0, 'Educational Subscription', 2],
    [ProductCode['Xmanager Enterprise'], 0x41, 0, 'Standard Subscription', 1],
    [ProductCode['Xmanager Enterprise'], 0x42, 0, 'Educational Subscription', 1],
    [ProductCode['Xmanager Enterprise'], 0x47, 0, 'Standard Subscription', 2],
    [ProductCode['Xmanager Enterprise'], 0x48, 0, 'Educational Subscription', 2],
    [ProductCode['Xmanager Enterprise'], 0x4B, 0, 'Standard Subscription', 1],
    [ProductCode['Xmanager Enterprise'], 0x4C, 0, 'Educational Subscription', 1],
    [ProductCode['Xmanager Enterprise'], 0x51, 0, '', 2],
    [ProductCode['Xmanager Enterprise'], 0x52, 0, 'Educational', 2],
    [ProductCode['Xmanager Enterprise'], 0x55, 0, '', 1],
    [ProductCode['Xmanager Enterprise'], 0x56, 0, 'Educational', 1],
    [ProductCode['Xmanager Enterprise'], 0x60, 0, 'Standard', 1],
    [ProductCode['Xmanager Enterprise'], 0x61, 0, 'Standard', 2],
    [ProductCode['Xmanager Enterprise'], 0x62, 0, 'Standard', 1],
    [ProductCode['Xmanager Enterprise'], 0x63, 0, 'Standard', 2],
]

ProductPublishList = (
    {'ProductName': 'Xmanager', 'Version': 2, 'PublishDate': datetime.date(2003, 1, 1)},
    {'ProductName': 'Xshell', 'Version': 2, 'PublishDate': datetime.date(2004, 10, 1)},

    {'ProductName': 'Xmanager', 'Version': 3, 'PublishDate': datetime.date(2007, 1, 1)},
    {'ProductName': 'Xshell', 'Version': 3, 'PublishDate': datetime.date(2007, 1, 1)},
    {'ProductName': 'Xlpd', 'Version': 3, 'PublishDate': datetime.date(2007, 1, 1)},
    {'ProductName': 'Xftp', 'Version': 3, 'PublishDate': datetime.date(2007, 1, 1)},
    {'ProductName': 'Xmanager Enterprise', 'Version': 3, 'PublishDate': datetime.date(2007, 1, 1)},

    {'ProductName': 'Xmanager', 'Version': 4, 'PublishDate': datetime.date(2010, 8, 1)},
    {'ProductName': 'Xshell', 'Version': 4, 'PublishDate': datetime.date(2010, 8, 1)},
    {'ProductName': 'Xlpd', 'Version': 4, 'PublishDate': datetime.date(2010, 8, 1)},
    {'ProductName': 'Xftp', 'Version': 4, 'PublishDate': datetime.date(2010, 8, 1)},
    {'ProductName': 'Xmanager Enterprise', 'Version': 4, 'PublishDate': datetime.date(2010, 8, 1)},
    {'ProductName': 'Xmanager', 'Version': 5, 'PublishDate': datetime.date(2014, 4, 28)},
    {'ProductName': 'Xshell', 'Version': 5, 'PublishDate': datetime.date(2014, 4, 28)},
    {'ProductName': 'Xlpd', 'Version': 5, 'PublishDate': datetime.date(2014, 4, 28)},
    {'ProductName': 'Xftp', 'Version': 5, 'PublishDate': datetime.date(2014, 4, 28)},
    {'ProductName': 'Xmanager Enterprise', 'Version': 5, 'PublishDate': datetime.date(2014, 4, 28)},
    {'ProductName': 'Xmanager', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)},
    {'ProductName': 'Xshell', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)},
    {'ProductName': 'Xshell Plus', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)},
    {'ProductName': 'Xlpd', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)},
    {'ProductName': 'Xftp', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)},
    {'ProductName': 'Xmanager Enterprise', 'Version': 6, 'PublishDate': datetime.date(2018, 4, 29)}
)


def get_check_sum(pre_product_key: str):
    check_sum = 1
    for i in range(0, len(pre_product_key)):
        if pre_product_key[i] != '-' and pre_product_key[i] != '8' and pre_product_key[i] != '9':
            place = int(pre_product_key[i])
            check_sum = (9 - place) * check_sum % -1000
    check_sum = (check_sum + int(pre_product_key[9])) % 1000
    return check_sum


def generate_product_key(issue_date: datetime.date,
                         product_name: str,
                         product_version: int,
                         number_license: int):
    if issue_date.year < 2002:
        raise ValueError('IssueDate cannot be earlier than 2002.')
    if issue_date > datetime.date.today() + datetime.timedelta(days=7):
        raise ValueError('IssueDate cannot be later than today after a week.')
    if number_license < 0 or number_license > 999:
        raise ValueError('NumberOfLicense must vary from 0 to 999.')

    for item in ProductPublishList:
        if item['ProductName'] == product_name and item['Version'] == product_version:
            if item['PublishDate'] > issue_date:
                raise ValueError('IssueDate cannot be earlier than the publish date.')
            break
        # if item['ProductName'] == ProductName and str(item['Version']) != ProductVersion:
        #     raise ValueError('Invalid product.')

    pre_product_key = '%02d%02d%02d-%02d%d%03d-%03d' % (issue_date.year - 2000,
                                                        issue_date.month,
                                                        issue_date.day,
                                                        0x0B,
                                                        ProductCode[product_name],
                                                        random.randint(0, 999),
                                                        number_license)
    check_sum = get_check_sum(pre_product_key)
    product_key = pre_product_key + '%03d' % check_sum
    return product_key


def generate_key(product_name: str, product_version: int):
    return generate_product_key(
        datetime.date(datetime.datetime.now().year,
                      datetime.datetime.now().month,
                      datetime.datetime.now().day), product_name,
        product_version, 999)


if __name__ == '__main__':
    msg = generate_key(sys.argv[1], int(sys.argv[2]))
    print(msg)
