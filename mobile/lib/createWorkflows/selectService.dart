import 'package:flutter/material.dart';

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
  String title = "Choose a service";
  int serviceIndex = 0;
  List<Widget> services = [];

  Future<void> redirectToSelectActionOrReactionPage(int servicesIndex) async {
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
        if (isAnAction) {
          Globaldata.actionServiceIndex = servicesIndex;
        } else {
          Globaldata.reactionServiceIndex = servicesIndex;
        }
        CreateData.selectOrUpdateWorkflow = false;
        Navigator.pop(Globaldata.myContext, CreateData.actualServiceList[servicesIndex].name);
      }
    }
  }

  Widget serviceTiles(int index) {
    if (index + 1 < CreateData.actualServiceList.length && index % 2 == 0) {
      serviceIndex++;
      return Row(
        children: [
          Padding(padding: EdgeInsets.only(right: 5)),
          SquareServicesCard(
            name: CreateData.actualServiceList[index].name,
            color: HexColor(CreateData.actualServiceList[index].color),
            imageUrl: CreateData.actualServiceList[index].icon,
            onTap: () {
              if (isAnAction) {
                redirectToSelectActionOrReactionPage(index);
              } else {
                redirectToSelectActionOrReactionPage(index);
              }
            },
          ),
          SquareServicesCard(
            name: CreateData.actualServiceList[index + 1].name,
            color: HexColor(CreateData.actualServiceList[index + 1].color),
            imageUrl: CreateData.actualServiceList[index + 1].icon,
            onTap: () {
              if (isAnAction) {
                redirectToSelectActionOrReactionPage(index + 1);
              } else {
                redirectToSelectActionOrReactionPage(index + 1);
              }
            },
          ),
        ],
      );
    }
    return Row(
      children: [
        Padding(padding: EdgeInsets.only(right: 5)),
        SquareServicesCard(
          name: CreateData.actualServiceList[index].name,
          color: HexColor(CreateData.actualServiceList[index].color),
          imageUrl: CreateData.actualServiceList[index].icon,
          onTap: () {
            if (isAnAction) {
              redirectToSelectActionOrReactionPage(index);
            } else {
              redirectToSelectActionOrReactionPage(index);
            }
          },
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
