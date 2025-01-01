import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;

import "../custom_widget.dart";
import "../globalData.dart";
import "selectActionOrReaction.dart";
import 'createWorkflowsData.dart';
import "createWidget.dart";

class SelectServicePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  SelectService createState() => SelectService();
}

class SelectService extends State<APICall> {
  String title = "Select action service";
  int serviceIndex = 0;
  List<Widget> services = [];

  Future<void> setActionOrReactionServiceList() async {
    if (isAnAction == true) {
      title = "Select action service";
      CreateData.serviceNameList = Globaldata.actionServicesName;
      CreateData.serviceColorList = Globaldata.actionServicesColor;
      CreateData.serviceIconsList = Globaldata.actionServicesIcons;
      return;
    }
    setState(() {
      title = "Select reaction service";
      CreateData.serviceNameList = Globaldata.reactionServicesName;
      CreateData.serviceColorList = Globaldata.reactionServicesColor;
      CreateData.serviceIconsList = Globaldata.reactionServicesIcons;
    });
  }

  Widget serviceTiles(int index) {
    if (index + 1 < CreateData.serviceNameList.length && index % 2 == 0) {
      serviceIndex++;
      return Row(
        children: [
          MySquareServicesCard(
            servicesIndex: index++,
          ),
          MySquareServicesCard(
            servicesIndex: index,
          ),
        ],
      );
    }
    return Row(
      children: [
        MySquareServicesCard(
          servicesIndex: index,
        ),
      ],
    );
  }

  @override
  void initState() {
    super.initState();
    setState(() {
      setActionOrReactionServiceList();
      for (serviceIndex = 0;
        serviceIndex < CreateData.serviceNameList.length;
        serviceIndex++) {
        services.add(serviceTiles(serviceIndex));
      }
    });
  }

  @override
  Widget build(BuildContext context) {
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

class MySquareServicesCard extends StatelessWidget {
  const MySquareServicesCard({
    super.key,
    required this.servicesIndex,
  });

  final int servicesIndex;

  Future<void> redirectToSelectActionOrReactionPage() async {
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(
          builder: (context) => SelectActionOrReactionPage(
            serviceIndex: servicesIndex,
          )),
    );
    if (isAnAction && CreateData.serviceAction.name.length != 0 || !isAnAction && CreateData.serviceReaction.name.length != 0) {
      Navigator.pop(Globaldata.myContext, CreateData.serviceNameList[servicesIndex]);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (servicesIndex > Globaldata.actionServicesName.length) {
      return MyDivider();
    }
    developer.log("https://mathisbrehm.fr/${CreateData.serviceIconsList[servicesIndex]}");
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          color: HexColor(CreateData.serviceColorList[servicesIndex]),
          border: Border.all(
            color: HexColor(CreateData.serviceColorList[servicesIndex]),
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
                padding:
                    const EdgeInsets.symmetric(horizontal: 5, vertical: 30),
                child: Column(children: [
                  SafeWebImage(
                      url:
                          "${Globaldata.domainName}${CreateData.serviceIconsList[servicesIndex]}"),
                  AutoSizeText(
                    CreateData.serviceNameList[servicesIndex],
                    style: const TextStyle(
                        fontSize: 25,
                        fontWeight: FontWeight.bold,
                        color: Colors.white),
                    maxLines: 1,
                  )
                ]),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
