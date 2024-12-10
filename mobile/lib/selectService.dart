import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;

import "custom_widget.dart";
import "globalData.dart";
import "selectActionOrReaction.dart";

class SelectServicePage extends StatelessWidget {
  const SelectServicePage({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(isAnAction: isAnAction),
    );
  }
}

class APICall extends StatefulWidget {
  const APICall({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  SelectService createState() => SelectService(isAnAction: isAnAction);
}

class SelectService extends State<APICall> {
  SelectService({
    required this.isAnAction,
  });
  final bool isAnAction;
  String title = "Select action service";
  List<String> serviceNameList = Globaldata.actionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons;
  int serviceIndex = 0;

  Future<void> setReactionService() async {
    setState(() {
      title = "Select reaction service";
      serviceNameList = Globaldata.reactionServicesName;
      serviceColorList = Globaldata.reactionServicesColor;
      serviceIconsList = Globaldata.reactionServicesIcons;
    });
  }

  Widget serviceTiles(int index) {
    if (index + 1 < serviceNameList.length && index % 2 == 0) {
      serviceIndex++;
      developer.log("index : $index, index2 : ${index + 1}");
      return Row(
        children: [
          MySquareServicesCard(
            servicesIndex: index++,
            isAnAction: isAnAction,
            serviceColorList: serviceColorList,
            serviceNameList: serviceNameList,
            serviceIconsList: serviceIconsList,
          ),
          MySquareServicesCard(
            servicesIndex: index,
            isAnAction: isAnAction,
            serviceColorList: serviceColorList,
            serviceNameList: serviceNameList,
            serviceIconsList: serviceIconsList,
          ),
        ],
      );
    }
    developer.log("index : $index");
    return Row(
      children: [
        MySquareServicesCard(
          servicesIndex: index,
          isAnAction: isAnAction,
          serviceColorList: serviceColorList,
          serviceNameList: serviceNameList,
          serviceIconsList: serviceIconsList,
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    if (isAnAction != true) {
      setReactionService();
    }
    List<Widget> services = [];
    for (serviceIndex = 0; serviceIndex < serviceNameList.length; serviceIndex++) {
      developer.log("index : ${serviceIndex}, list size : ${serviceNameList.length}");
      services.add(
        serviceTiles(serviceIndex)
      );
    }

    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectServiceTopBar(title: title),

            MyDivider(),

            Column(
              children: services,
            ),
          ],
        ),
      ),
    );
  }
}

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
  });

  final String title;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {Navigator.pop(Globaldata.myContext);},
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(title, style: TextStyle( fontSize: 30, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),
        ),
      ],
    );
  }
}

class MySquareServicesCard extends StatelessWidget {
  const MySquareServicesCard({
    super.key,
    required this.servicesIndex,
    this.isAnAction = true,
    required this.serviceNameList,
    required this.serviceColorList,
    required this.serviceIconsList,
  });

  final int servicesIndex;
  final bool isAnAction;
  final List<String> serviceNameList;
  final List<String> serviceColorList;
  final List<String> serviceIconsList;

  Future<void> redirectToSelectActionOrReactionPage() async {
      await Navigator.push(Globaldata.myContext,MaterialPageRoute(builder: (context) => SelectActionOrReactionPage(isAnAction: isAnAction, serviceIndex: servicesIndex,)),);
      Navigator.pop(Globaldata.myContext, serviceNameList[servicesIndex]);
  }

  @override
  Widget build(BuildContext context) {
    if (servicesIndex >= Globaldata.actionServicesName.length) {
      return MyDivider();
    }
    developer.log("https://mathisbrehm.fr/${serviceIconsList[servicesIndex]}");
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          color: HexColor(serviceColorList[servicesIndex]),
          border: Border.all(
            color: HexColor(serviceColorList[servicesIndex]),
            width: 2.0,
          ),
          borderRadius: BorderRadius.circular(4.0),
        ),

        child: InkWell(
          onTap: () {
            if (isAnAction) {
              Globaldata.firstService = servicesIndex;
              redirectToSelectActionOrReactionPage();
            } else {
              Globaldata.secondService = servicesIndex;
              redirectToSelectActionOrReactionPage();
            }
          },
          child: Column(
            children: [
              SizedBox(
                width: MediaQuery.sizeOf(Globaldata.myContext).width / 2.5,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 30),
                child: Column(
                  children: [
                    SafeWebImage(
                        url: "${Globaldata.domainName}${serviceIconsList[servicesIndex]}"
                      ),
                    AutoSizeText(serviceNameList[servicesIndex], style: const TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,)
                  ]
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class HexColor extends Color {
  static int _getColorFromHex(String hexColor) {
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return int.parse(hexColor, radix: 16);
  }

  HexColor(final String hexColor) : super(_getColorFromHex(hexColor));
}