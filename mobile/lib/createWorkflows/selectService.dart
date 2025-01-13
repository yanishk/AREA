import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../custom_widget.dart";
import "../globalData.dart";
import "selectActionsOrReactions/selectActionOrReaction.dart";
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
      return;
    }
    setState(() {
      title = "Select reaction service";
    });
  }

  Widget serviceTiles(int index) {
    if (index + 1 < CreateData.actualServiceList.length && index % 2 == 0) {
      serviceIndex++;
      return Row(
        children: [
          Padding(padding: EdgeInsets.only(right: 5)),
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
        Padding(padding: EdgeInsets.only(right: 5)),
        MySquareServicesCard(
          servicesIndex: index,
        ),
      ],
    );
  }

  @override
  void initState() {
    super.initState();
    CreateData.actualServiceList = [];
    if (isAnAction) {
      CreateData.actionServiceList = [];
    } else {
      CreateData.reactionServiceList = [];
    }
    setState(() {
      setActionOrReactionServiceList();
      for (int i = 0; i < Globaldata.serviceList.length; i++) {
        if (isAnAction) {
          if (Globaldata.serviceList[i].hasActions) {
            CreateData.actualServiceList.add(Globaldata.serviceList[i]);
            CreateData.actionServiceList.add(Globaldata.serviceList[i]);
          }
        } else if (Globaldata.serviceList[i].hasReactions) {
          CreateData.actualServiceList.add(Globaldata.serviceList[i]);
          CreateData.reactionServiceList.add(Globaldata.serviceList[i]);
        }
      }
      for (serviceIndex = 0;
        serviceIndex < CreateData.actualServiceList.length;
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
        )
      ),
    );
    if (isAnAction && CreateData.serviceAction.name.length != 0 || !isAnAction && CreateData.serviceReaction.name.length != 0) {
      if (CreateData.selectOrUpdateWorkflow) {
        CreateData.selectOrUpdateWorkflow = false;
        Navigator.pop(Globaldata.myContext, CreateData.actualServiceList[servicesIndex].name);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (servicesIndex > CreateData.actualServiceList.length) {
      return Column();
    }
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
      child: Container(
        width: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        height: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        decoration: BoxDecoration(
          color: HexColor(CreateData.actualServiceList[servicesIndex].color),
          border: Border.all(
            color: HexColor(CreateData.actualServiceList[servicesIndex].color),
            width: 2.0,
          ),
          borderRadius: BorderRadius.circular(6.0),
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
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SafeWebImage(
                url: "${Globaldata.domainName}${CreateData.actualServiceList[servicesIndex].icon}",
                width: 50,
                height: 50,
              ),
              Padding(padding: EdgeInsets.only(top: 20)),
              AutoSizeText(
                CreateData.actualServiceList[servicesIndex].name,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white
                ),
                maxLines: 1,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
